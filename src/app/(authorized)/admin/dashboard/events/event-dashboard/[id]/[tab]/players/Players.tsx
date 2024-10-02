import {
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
    // const event = await getEventById(id);

    const poolPlayerInterests = await generatePoolPlayerInterestById('1AJUD3MiLuvnhf5AcXHD');
    console.log('poolPlayerInterests', poolPlayerInterests);

    // const dummyPlayerInrerestData: PlayerInterest = {
    //     interestLevel: InterestLevel.VeryInterested,
    //     poolEventId: '',
    //     participantId: '',
    //     firstName: 'Kari',
    //     lastName: 'Nordmann',
    //     isOver18: true,
    //     ticketCategoryID: 0,
    //     ticketCategory: 'Festivalpass Ungdom/student (13-25/30år) Early-bird',
    //     conPlayers: [
    //         {
    //             id: undefined,
    //             participantId: '',
    //             firstName: '',
    //             lastName: '',
    //             interestLevel: InterestLevel.VeryInterested,
    //             poolEventId: '',
    //             poolEventTitle: 'Et gøyalt spill',
    //             poolName: PoolName.fridayEvening,
    //             isFirstChoice: true,
    //             isGameMaster: false,
    //             createdAt: '',
    //             createdBy: '',
    //             updateAt: '',
    //             updatedBy: '',
    //         },
    //         {
    //             id: undefined,
    //             participantId: '',
    //             firstName: '',
    //             lastName: '',
    //             interestLevel: InterestLevel.SomewhatInterested,
    //             poolEventId: '',
    //             poolEventTitle: '"Hva f**n skjedde!?!??" -- ghost/echo av John Harper',
    //             poolName: PoolName.saturdayMorning,
    //             isFirstChoice: false,
    //             isGameMaster: true,
    //             createdAt: '',
    //             createdBy: '',
    //             updateAt: '',
    //             updatedBy: '',
    //         },
    //         {
    //             id: undefined,
    //             participantId: '',
    //             firstName: '',
    //             lastName: '',
    //             interestLevel: InterestLevel.Interested,
    //             poolEventId: '',
    //             poolEventTitle: 'Random dnd event',
    //             poolName: PoolName.saturdayMorning,
    //             isFirstChoice: false,
    //             isGameMaster: true,
    //             createdAt: '',
    //             createdBy: '',
    //             updateAt: '',
    //             updatedBy: '',
    //         },

    return (
        <Paper sx={{ padding: '1rem' }}>
            <Typography variant="h1">Spillere:</Typography>
            <Card sx={{ backgroundColor: 'rgb(55, 59, 87)' }}>
                <CardHeader title="Prioriterg" />
                <CardActions>
                    <FormGroup sx={{ display: 'flex', flexDirection: 'row', gap: '1rem' }}>
                        <FormControlLabel control={<Switch />} label="Under 18" labelPlacement="start" />
                        <FormControlLabel control={<Switch />} label="Over 18" labelPlacement="start" />
                        <FormControlLabel control={<Switch />} label="Søndag Dagspass Barn" labelPlacement="start" />
                    </FormGroup>
                </CardActions>
            </Card>
            <Paper sx={{ backgroundColor: 'rgpa(0,0,0,0.1)', padding: '1rem' }}>
                <FormControl variant="standard">
                    <InputLabel htmlFor="input-with-icon-adornment">Søk etter deltager</InputLabel>
                    <Input
                        id="input-with-icon-adornment"
                        endAdornment={
                            <InputAdornment position="start">
                                <SearchIcon />
                            </InputAdornment>
                        }
                    />
                </FormControl>
            </Paper>
            <Typography variant="h2">Spillere:</Typography>
            {poolPlayerInterests.map((poolPlayerInterest) => (
                <PlayerInterestInfo poolName={PoolName.saturdayEvening} playerInterest={poolPlayerInterest} />
            ))}
        </Paper>
    );
};

export default Players;
