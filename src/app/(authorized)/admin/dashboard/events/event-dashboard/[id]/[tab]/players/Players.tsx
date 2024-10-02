import {
    Box,
    Card,
    CardActions,
    CardHeader,
    FormControl,
    FormControlLabel,
    FormGroup,
    Input,
    InputAdornment,
    InputLabel,
    Paper,
    Switch,
    Typography,
} from '@mui/material';
import PlayerInterestInfo from './components/PlayerInterestInfo';
import { PlayerInterest } from '$lib/types';
import { InterestLevel, PoolName } from '$lib/enums';
import SearchIcon from '@mui/icons-material/Search';
import { getEventById, getPoolEventById } from '$app/(public)/components/lib/serverAction';
import { generatePoolPlayerInterestById } from './components/lib/actions';

type Props = {
    id: string;
};

const Players = async ({ id }: Props) => {
    const event = await getEventById(id);

    return (
        <Box>
            <Typography variant="h1">Spillere:</Typography>
            <Paper sx={{ padding: '1rem' }}></Paper>
        </Box>
    );
};

export default Players;
