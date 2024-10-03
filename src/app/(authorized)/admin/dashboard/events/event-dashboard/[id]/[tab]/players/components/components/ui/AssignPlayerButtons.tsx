'use client';
import { CircularProgress, FormControlLabel, FormGroup, Switch, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';
import { assignPlayer } from '../lib/actions';
import { useTransition } from 'react';

type Props = {
    participantId: string;
    poolEventId: string;
    isAssigned: boolean;
    isGameMaster: boolean;
    poolPlayerId?: string;
};

const AssignPlayerButtons = ({ poolPlayerId, participantId, poolEventId, isAssigned, isGameMaster }: Props) => {
    const router = useRouter();
    const [isPending, startTransition] = useTransition();
    const handleChanges = async (event: React.ChangeEvent<HTMLInputElement>) => {
        startTransition(async () => {
            const { checked } = event.target;
            const name = event.target.name as 'isAssigned' | 'isGameMaster';
            const isAssigned = name === 'isAssigned' ? checked : false;
            const isGameMaster = name === 'isGameMaster' ? checked : false;
            await assignPlayer(participantId, poolEventId, isAssigned, isGameMaster, poolPlayerId);
            console.log('event', event, 'name', name, 'checked', checked);
            router.refresh();
        });
    };

    return isPending ?
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
        :   <>
                <Typography sx={{ display: 'inline-block' }}>Oppdaterer listen med tildelte</Typography>{' '}
                <CircularProgress size="1.5rem" />
            </>;
};
export default AssignPlayerButtons;
