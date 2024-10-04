import { InterestLevel, PoolName, RoomName } from '$lib/enums';
import { PlayerInterest } from '$lib/types';
import { Typography } from '@mui/material';
import PlayerInterestInfo from './PlayerInterestInfo';

type Props = {
    assignedPlayers: PlayerInterest[];
    poolName: PoolName;
};

const AssignedPlayers = async ({ assignedPlayers, poolName }: Props) => {
    return (
        <>
            <Typography
                variant="h1"
                sx={{ scrollMarginTop: 'calc(var(--app-bar-height-desktop, 0px) + 146px)' }}
                id="assigned-players"
            >
                Assigned Players
            </Typography>
            {assignedPlayers.map((participantInterest) => (
                <PlayerInterestInfo
                    key={participantInterest.participantId}
                    playerInterest={participantInterest}
                    poolName={poolName}
                />
            ))}
        </>
    );
};

export default AssignedPlayers;
