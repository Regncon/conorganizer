import { PlayerInterest } from '$lib/types';
import { Box, Divider, FormControlLabel, FormGroup, Paper, Switch, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from '../../../lib/helpers';
import PoolPlayer from './components/PreviousGamesPlayed';
import { PoolName } from '$lib/enums';
import Over18 from '$ui/participant/Over18';
import ParticipantAvatar from '$ui/participant/ParticipantAvatar';
import AssignPlayerButtons from './ui/AssignPlayerButtons';

type Props = {
    poolName: PoolName;
    playerInterest: PlayerInterest;
};

const PlayerInterestInfo = ({ poolName, playerInterest }: Props) => {
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
                        />
                        <Over18 over18={playerInterest.isOver18} />
                        <Typography component={'i'}>{playerInterest.ticketCategory}</Typography>
                    </Box>
                </Box>
            </Box>
            <Box>
                <PoolPlayer
                    currentPoolName={poolName}
                    previousPoolName={PoolName.fridayEvening}
                    conPlayers={playerInterest.conPlayers}
                />
                <PoolPlayer
                    currentPoolName={poolName}
                    previousPoolName={PoolName.saturdayMorning}
                    conPlayers={playerInterest.conPlayers}
                />
                <PoolPlayer
                    currentPoolName={poolName}
                    previousPoolName={PoolName.saturdayEvening}
                    conPlayers={playerInterest.conPlayers}
                />
                <PoolPlayer
                    currentPoolName={poolName}
                    previousPoolName={PoolName.sundayMorning}
                    conPlayers={playerInterest.conPlayers}
                />
            </Box>
        </Paper>
    );
};

export default PlayerInterestInfo;
