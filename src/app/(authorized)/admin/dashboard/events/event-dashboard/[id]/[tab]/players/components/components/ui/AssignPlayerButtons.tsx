'use client';
import { FormControlLabel, FormGroup, Switch } from '@mui/material';
import { useRouter } from 'next/navigation';

type Props = {
    participantId: string;
    poolEventId: string;
    isAssigned: boolean;
    isGameMaster: boolean;
};

const AssignPlayerButtons = ({ participantId, poolEventId, isAssigned, isGameMaster }: Props) => {
    const router = useRouter();
    const handleChanges = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const { name, checked } = event.target;
        console.log('event', event, 'name', name, 'checked', checked);
        router.refresh();
    };

    return (
        <FormGroup sx={{ display: 'flex', flexDirection: 'row', gap: '1rem' }}>
            <FormControlLabel
                control={<Switch defaultChecked={isAssigned} name="isAssigned" onChange={handleChanges} />}
                label="Tildel plass"
                labelPlacement="start"
            />
            <FormControlLabel
                control={<Switch defaultChecked={isGameMaster} name="isGameMaster" onChange={handleChanges} />}
                label="GM"
                labelPlacement="start"
            />
        </FormGroup>
    );
};
export default AssignPlayerButtons;
