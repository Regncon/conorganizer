import { PlayerInterest } from '$lib/types';
import { Box, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from '../../../lib/helpers';

import { PoolName } from '$lib/enums';
import Over18 from '$ui/participant/Over18';
import ParticipantAvatar from '$ui/participant/ParticipantAvatar';
import AssignPlayerButtons from './ui/AssignPlayerButtons';
import { generatePlayerInPoolMap } from './lib/playerInPoolHelper';
import PreviousGamesPlayed from './components/PreviousGamesPlayed';

type Props = {
    poolName: PoolName;
    playerInterest: PlayerInterest;
};

const PlayerInterestInfo = ({ poolName, playerInterest }: Props) => {
    console.log('playerInterest', playerInterest.playerInPools);
    let playersInPool = generatePlayerInPoolMap(playerInterest.playerInPools);

    return (
        <Paper elevation={2} sx={{ marginTop: '1rem', padding: '1rem' }}>
            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'flex-start',
                    gap: '1rem',
                    alignItems: 'center',
                    flexDirection: 'row',
                }}
            >
                <Box>
                    <Image
                        src={interestLevelToImage[playerInterest.interestLevel]}
                        alt={InterestLevelToLabel[playerInterest.interestLevel]}
                        width={100}
                        height={60}
                    />
                </Box>

                <Box>
                    <Box sx={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: '1rem' }}>
                        <Typography>{InterestLevelToLabel[playerInterest.interestLevel]}</Typography>
                        <ParticipantAvatar
                            firstName={playerInterest.firstName}
                            lastName={playerInterest.lastName}
                            header
                        />
                        <AssignPlayerButtons
                            participantId={playerInterest.participantId}
                            poolEventId={playerInterest.poolEventId}
                            isAssigned={playerInterest.isAssigned}
                            isGameMaster={playerInterest.isGameMaster}
                            poolPlayerId={}
                        />
                        <Over18 over18={playerInterest.isOver18} />
                        <Typography component={'i'}>{playerInterest.ticketCategory}</Typography>
                    </Box>
                </Box>
            </Box>
            <Box>
                {[...playersInPool.entries()].map(([day, poolPlayer]) => {
                    const hasNoPlayerOnDay = Object.keys(poolPlayer).length === 0;
                    return (
                        <PreviousGamesPlayed
                            poolName={day}
                            poolPlayer={poolPlayer}
                            hasNoPlayerOnDay={hasNoPlayerOnDay}
                        />
                    );
                })}
            </Box>
        </Paper>
    );
};

export default PlayerInterestInfo;
