import { PoolPlayer } from '$lib/types';
import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from '../../../../lib/helpers';
import { getTranslatedDay } from '$app/(public)/components/lib/helpers/translation';
import { PoolName } from '$lib/enums';
import WarningIcon from '@mui/icons-material/Warning';
import GamemasterIcon from '$lib/components/icons/GameMasterIcon';

type Props = {
    poolName: PoolName;
    poolPlayer: PoolPlayer;
    hasNoPlayerOnDay: boolean;
};
const PreviousGamesPlayed = ({ poolName, poolPlayer, hasNoPlayerOnDay }: Props) => {
    if (hasNoPlayerOnDay) {
        return (
            <Box>
                <Typography sx={{ fontWeight: 'bold' }}>{getTranslatedDay(poolName)}: - </Typography>
            </Box>
        );
    }
    return (
        <>
            <Box
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
                <Typography sx={{ fontWeight: 'bold' }}>{getTranslatedDay(poolName)}: </Typography>
                <Image
                    src={interestLevelToImage[poolPlayer.interestLevel]}
                    alt={InterestLevelToLabel[poolPlayer.interestLevel]}
                    width={50}
                    height={25}
                />
                <Typography>{InterestLevelToLabel[poolPlayer.interestLevel]}</Typography>
                <Box sx={{ display: 'flex', flexDirection: 'row', gap: '1rem' }}>
                    <Typography component={'i'}> {poolPlayer.poolEventTitle}</Typography>
                    {poolPlayer.isFirstChoice && (
                        <Box sx={{ display: 'flex', gap: '1rem', color: 'warning.main' }}>
                            <WarningIcon />
                            <Typography component={'i'}>Fått førstevalg</Typography>
                        </Box>
                    )}
                    {poolPlayer.isGameMaster && (
                        <Box sx={{ display: 'flex', gap: '1rem', color: 'success.main' }}>
                            <GamemasterIcon color="success" />
                            <Typography>Spilleder! :D</Typography>
                        </Box>
                    )}
                </Box>
            </Box>
        </>
    );
};
export default PreviousGamesPlayed;
