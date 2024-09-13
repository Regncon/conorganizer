import { Box, Stack, Typography } from '@mui/material';
import WarningIcon from '@mui/icons-material/Warning';
import { PoolName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { poolTitles } from '../../lib/helpers';

type props = {
    poolName: PoolName;
    conEvent: ConEvent;
};

const UnwantedPoolByGm = ({ poolName, conEvent }: props) => {
    let showUnwanted = false;
    switch (poolName) {
        case PoolName.fridayEvening:
            showUnwanted = conEvent.unwantedFridayEvening;
            break;
        case PoolName.saturdayMorning:
            showUnwanted = conEvent.unwantedSaturdayMorning;
            break;
        case PoolName.saturdayEvening:
            showUnwanted = conEvent.unwantedSaturdayEvening;
            break;
        case PoolName.sundayMorning:
            showUnwanted = conEvent.unwantedSundayMorning;
            break;
    }
    if (!showUnwanted) {
        return null;
    }
    return (
        <Stack direction="row">
            <Box sx={{ display: 'inherit', color: 'warning.main' }}>
                <WarningIcon />
                <Typography component={'i'}>Gm ønsker ikke {poolTitles[poolName]}</Typography>
            </Box>
        </Stack>
    );
};

export default UnwantedPoolByGm;
