'use client';
import { FormControlLabel, FormGroup, Switch } from '@mui/material';
import { useRouter } from 'next/navigation';
import { assignPlayer } from '../lib/actions';

type Props = {
    poolPlayerId: string;
    participantId: string;
    poolEventId: string;
    isAssigned: boolean;
    isGameMaster: boolean;
};

const AssignPlayerButtons = ({ poolPlayerId, participantId, poolEventId, isAssigned, isGameMaster }: Props) => {
    const router = useRouter();
    const handleChanges = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const { name, checked } = event.target;
        const isAssigned = name === 'isAssigned' ? checked : false;
        const isGameMaster = name === 'isGameMaster' ? checked : false;
        await assignPlayer(poolPlayerId, participantId, poolEventId, isAssigned, isGameMaster);
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
