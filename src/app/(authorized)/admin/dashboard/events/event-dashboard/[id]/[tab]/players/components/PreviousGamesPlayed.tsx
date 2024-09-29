import { ConPlayer, Interest, Participant, PlayerInterest } from '$lib/types';
import { Box, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from '../../lib/helpers';
import { getTranslatedDay } from '$app/(public)/components/lib/helpers/translation';
import { PoolName } from '$lib/enums';
import WarningIcon from '@mui/icons-material/Warning';
import GamemasterIcon from '$lib/components/icons/GameMasterIcon';

type Props = {
    poolName: PoolName;
    conPlayers: ConPlayer[];
};
const PreviousGamesPlayed = ({ poolName, conPlayers }: Props) => {
    console.log('PreviousGamesPlayed', poolName, conPlayers);
    const filteredConPlayers = conPlayers.filter((conPlayer) => conPlayer.poolName === poolName);

    if (filteredConPlayers === undefined || filteredConPlayers.length === 0) {
        return (
            <Box>
                <Typography sx={{ fontWeight: 'bold' }}>{getTranslatedDay(poolName)}: - </Typography>
            </Box>
        );
    }
    const playersInPool = filteredConPlayers.filter((conPlayer) => conPlayer.poolName === poolName);
    return (
        <Box sx={{ backgroundColor: 'rgba(0,0,0,0.1)', marginBlock: '1rem' }}>
            <Typography sx={{ fontWeight: 'bold' }}>{getTranslatedDay(poolName)}: </Typography>
            {playersInPool.map((conPlayer, index) => (
                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: 'column',
                        backgroundColor: 'rgba(0,0,0,0.2)',
                        marginBlock: '1rem',
                        gap: '0.5rem',
                    }}
                >
                    <Typography component={'i'}> {conPlayer.poolEventTitle}</Typography>
                    <Box key={index} sx={{ display: 'flex', gap: '1rem' }}>
                        <Image
                            src={interestLevelToImage[conPlayer.interestLevel]}
                            alt={InterestLevelToLabel[conPlayer.interestLevel]}
                            width={50}
                            height={25}
                        />
                        <Typography>{InterestLevelToLabel[conPlayer.interestLevel]}</Typography>
                    </Box>
                    {conPlayer.isFirstChoice && (
                        <Box sx={{ display: 'flex', gap: '1rem', color: 'warning.main' }}>
                            <WarningIcon />
                            <Typography component={'i'}>Fått førstevalg</Typography>
                        </Box>
                    )}
                    {conPlayer.isGameMaster && (
                        <Box sx={{ display: 'flex', gap: '1rem', color: 'success.main' }}>
                            <GamemasterIcon color="success" />
                            <Typography>Spilleder! :D</Typography>
                        </Box>
                    )}
                </Box>
            ))}
        </Box>
    );
};
export default PreviousGamesPlayed;
