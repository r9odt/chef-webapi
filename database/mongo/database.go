package mongo

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/JIexa24/chef-webapi/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConnector describes the structure of the connector and
// implements the Database interface.
type DBConnector struct {
	DSN          string
	DatabaseName string
	DB           *mongo.Client
	Logger       logging.Logger
}

// New return a new connector.
func New(dsn, databaseName string, l logging.Logger) (*DBConnector, error) {
	connector := &DBConnector{
		DatabaseName: databaseName,
		DSN:          dsn,
		Logger:       l,
	}
	var ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(dsn)
	db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	connector.DB = db
	l.Info("Connected to MongoDB!")
	return connector, nil
}

// Close connector.
func (db *DBConnector) Close() {
	if err := db.DB.Disconnect(context.TODO()); err != nil {
		db.Logger.Errorf("Close [db.DB.Disconnect]: %s", err.Error())
	}
}

// UploadFile uploads file to database.
func (db *DBConnector) UploadFile(path, filename string) error {
	var fileContent []byte
	var bucket *gridfs.Bucket
	bucket, err := gridfs.NewBucket(db.DB.Database(db.DatabaseName))
	if err != nil {
		return err
	}

	// Specify the Metadata option to include a "metadata" field in the files collection document.
	uploadOpts := options.GridFSUpload().SetMetadata(bson.D{{Key: "metadata tag", Value: filename}})
	uploadStream, err := bucket.OpenUploadStream(filename, uploadOpts)
	if err != nil {
		return err
	}
	fileContent, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	// Use SetWriteDeadline to force a timeout if the upload does not succeed in 2 seconds.
	if err = uploadStream.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return err
	}

	if _, err = uploadStream.Write(fileContent); err != nil {
		return err
	}
	if err = uploadStream.Close(); err != nil {
		return err
	}
	return nil
}

const collectionFilesName string = "fs.files"

// CheckFile check file in database.
func (db *DBConnector) CheckFile(filename string) bool {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionFilesName)
	filter := bson.D{primitive.E{Key: "filename", Value: filename}}
	result := collection.FindOne(context.TODO(), filter)
	return result.Err() != mongo.ErrNoDocuments
}

// DownloadFile downloads file from database.
func (db *DBConnector) DownloadFile(filename string) []byte {
	var fileContent = make([]byte, 0)
	var bucket *gridfs.Bucket
	bucket, err := gridfs.NewBucket(db.DB.Database(db.DatabaseName))
	if err != nil {
		db.Logger.Errorf("DownloadFile [gridfs.NewBucket]: %s", err.Error())
		return append(fileContent, []byte(err.Error())...)
	}

	downloadStream, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		// db.Logger.Errorf("DownloadFile [bucket.OpenDownloadStreamByName]: %s",
		// 	err.Error())
		return append(fileContent, []byte(err.Error())...)
	}

	// Use SetReadDeadline to force a timeout if the download does not succeed in 2 seconds.
	if err = downloadStream.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		db.Logger.Errorf("DownloadFile [downloadStream.SetReadDeadline]: %s",
			err.Error())
		return append(fileContent, []byte(err.Error())...)
	}

	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, downloadStream); err != nil {
		db.Logger.Errorf("DownloadFile [io.Copy]: %s", err.Error())
		return append(fileContent, []byte(err.Error())...)
	}
	if err := downloadStream.Close(); err != nil {
		db.Logger.Errorf("DownloadFile [downloadStream.Close]: %s", err.Error())
	}
	fileContent = fileBuffer.Bytes()
	return fileContent
}
