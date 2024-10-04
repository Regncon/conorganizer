import { InterestLevel, PoolName, RoomName } from '$lib/enums';
import { PlayerInterest } from '$lib/types';
import PlayerInterestInfo from './PlayerInterestInfo';

type Props = {
    assignedPlayers: PlayerInterest[];
    poolName: PoolName;
};

const AssignedPlayers = async ({ assignedPlayers, poolName }: Props) => {
    return (
        <>
            <h1>Assigned Players</h1>
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
