import { ConPlayer } from '$lib/types';
import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from '../../../../lib/helpers';
import { getTranslatedDay } from '$app/(public)/components/lib/helpers/translation';
import { PoolName } from '$lib/enums';
import WarningIcon from '@mui/icons-material/Warning';
import GamemasterIcon from '$lib/components/icons/GameMasterIcon';

type Props = {
    previousPoolName: PoolName;
    currentPoolName: PoolName;
    conPlayers: ConPlayer[];
};
const PreviousGamesPlayed = ({ previousPoolName, currentPoolName, conPlayers }: Props) => {
    // console.log('PreviousGamesPlayed', poolName, conPlayers);
    const filteredConPlayers = conPlayers.filter((conPlayer) => conPlayer.poolName === previousPoolName);

    // ckeck if the current pool is earlier than the previous pool
    if (previousPoolName < currentPoolName && (filteredConPlayers === undefined || filteredConPlayers.length === 0)) {
        return (
            <Box>
                <Typography sx={{ fontWeight: 'bold' }}>{getTranslatedDay(previousPoolName)}: - </Typography>
            </Box>
        );
    }
    const playersInPool = filteredConPlayers.filter((conPlayer) => conPlayer.poolName === previousPoolName);
    return (
        <>
            {playersInPool.map((conPlayer, index) => (
                <Box
                    key={index}
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: ' 9rem 3rem 10rem 1fr',
                        gap: '1rem',
                        alignItems: 'center',
                        backgroundColor: 'rgba(0,0,0,0.1)',
                        marginBlock: '0.2rem',
                        paddingBottom: '0.1rem',
                    }}
                >
                    <Typography sx={{ fontWeight: 'bold' }}>{getTranslatedDay(previousPoolName)}: </Typography>
                    <Image
                        src={interestLevelToImage[conPlayer.interestLevel]}
                        alt={InterestLevelToLabel[conPlayer.interestLevel]}
                        width={50}
                        height={25}
                    />
                    <Typography>{InterestLevelToLabel[conPlayer.interestLevel]}</Typography>
                    <Box sx={{ display: 'flex', flexDirection: 'row', gap: '1rem' }}>
                        <Typography component={'i'}> {conPlayer.poolEventTitle}</Typography>
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
                </Box>
            ))}
        </>
    );
};
export default PreviousGamesPlayed;
