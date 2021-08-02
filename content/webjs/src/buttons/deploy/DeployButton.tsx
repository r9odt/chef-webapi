import { FC,  ReactElement } from 'react';
import SendIcon from '@material-ui/icons/Send';
import Button, { ButtonProps } from '@material-ui/core/Button';
import {
    useUpdate,
    useNotify,
    useRefresh,
    Record,
} from 'react-admin';
import { deployersResource } from "../../App.js";

interface Props {
    deployAPIPath?: string;
    onlyResource?: boolean;
    selectedResource?: boolean;
    displayLabel?: string;
    deployResource: string;
    icon?: ReactElement;
    record?: Record;
    color?: string;
}

export type DeployButtonProps = Props & ButtonProps;

export const DeployIcon = <SendIcon />;
const DeployButton: FC<DeployButtonProps> = props => {
    const {
        deployAPIPath = deployersResource,
        deployResource = 'GuaranteedAbsolutelyDefinitelyNonExistentResource',
        selectedResource = false,
        onlyResource = false,
        displayLabel = 'Deploy',
        record = { id: 'GuaranteedAbsolutelyDefinitelyNonExistentID' },
        color = 'primary',
    } = props;
    const [update] = useUpdate();
    const notify = useNotify();
    const refresh = useRefresh();

    const action = (e: any) => {
        const resource = deployResource;
        const id = record.id;
        
        const resourcesListStub = ""
        const payload = { id, onlyResource, resourcesListStub, selectedResource };
        e.stopPropagation()
        update(deployAPIPath, resource, payload, {}, {
            onSuccess: () => {
                notify(
                    `Task for ${resource}: ${id} created`,
                    'info'
                );
                refresh();
            },
            onFailure: () => {
                notify(
                    `Task for ${resource}: ${id} not created`,
                    'warning'
                );
            },
        });
    }

    const displayedLabel = displayLabel
    return (
        <Button
            size='small'
            color={color}
            onClick={action}
            endIcon={<SendIcon />}
            type={'button'}
            aria-label={displayedLabel}
        >
            {displayedLabel}
        </Button>
    );
};

export default DeployButton;
