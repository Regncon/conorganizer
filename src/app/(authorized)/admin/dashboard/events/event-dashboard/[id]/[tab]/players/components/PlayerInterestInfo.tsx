import { Interest, Participant, PlayerInterest } from '$lib/types';
import { Box, Divider, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from '../../lib/helpers';
import { getTranslatedDay } from '$app/(public)/components/lib/helpers/translation';
import PreviousGamesPlayed from './PreviousGamesPlayed';
import { PoolName } from '$lib/enums';
import Over18 from '$ui/participant/Over18';
import ParticipantAvatar from '$ui/participant/ParticipantAvatar';
import { Padding } from '@mui/icons-material';

type Props = {
    playerInterest: PlayerInterest;
};

const PlayerInterestInfo = ({ playerInterest }: Props) => {
    return (
        <Paper elevation={2} sx={{ marginTop: '1rem', padding: '1rem' }}>
            <ParticipantAvatar firstName={playerInterest.firstName} lastName={playerInterest.lastName} header />
            <Box sx={{ display: 'flex', gap: '1rem' }}>
                <Image
                    src={interestLevelToImage[playerInterest.interestLevel]}
                    alt={InterestLevelToLabel[playerInterest.interestLevel]}
                    width={100}
                    height={60}
                />
                <Box>
                    <Typography>{InterestLevelToLabel[playerInterest.interestLevel]}</Typography>
                    <Over18 over18={playerInterest.isOver18} />
                </Box>
            </Box>
            <Typography component={'i'}>{playerInterest.ticketCategory}</Typography>
            <Divider sx={{ paddingBottom: '1rem' }} />

            <Box>
                <PreviousGamesPlayed poolName={PoolName.fridayEvening} conPlayers={playerInterest.conPlayers} />
                <PreviousGamesPlayed poolName={PoolName.saturdayMorning} conPlayers={playerInterest.conPlayers} />
                <PreviousGamesPlayed poolName={PoolName.saturdayEvening} conPlayers={playerInterest.conPlayers} />
                <PreviousGamesPlayed poolName={PoolName.sundayMorning} conPlayers={playerInterest.conPlayers} />
            </Box>
        </Paper>
    );
};
export default PlayerInterestInfo;
