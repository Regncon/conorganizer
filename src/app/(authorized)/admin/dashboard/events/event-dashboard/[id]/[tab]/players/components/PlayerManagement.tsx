import {
    Box,
    Button,
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
    TextField,
    Typography,
} from '@mui/material';
import { PoolName } from '$lib/enums';
import SearchIcon from '@mui/icons-material/Search';
import { generatePoolPlayerInterestById } from './lib/actions';
import PlayerInterestInfo from './components/PlayerInterestInfo';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';
import AssigendPlayers from './components/AssigendPlayers';

type Props = {
    id: string | undefined;
    poolName: PoolName;
    maxNumberOfPlayers: number;
};

const PlayerManagement = async ({ id, poolName, maxNumberOfPlayers }: Props) => {
    if (!id) {
        return <Typography variant="h1">Event ikke satt opp i denne puljen</Typography>;
    }
    // const event = await getEventById(id);

    //TODO:Må ha puljenavn og hente data relatert til pulje du er i
    // const poolPlayerInterests = await generatePoolPlayerInterestById(id, poolName);
    const poolPlayerInterests = await generatePoolPlayerInterestById(id);
    // console.log('poolPlayerInterests', poolPlayerInterests);

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
            <FormGroup sx={{ display: 'flex', flexDirection: 'row', gap: '1rem', placeItems: 'center' }}>
                <Typography variant="h1">{translatedDays.get(poolName)}</Typography>
                <Button sx={{ maxWidth: 'fit-content', maxHeight: 'fit-content' }} variant="contained" color="primary">
                    Algotitme
                </Button>
                <FormControlLabel control={<Switch />} label="Under 18" labelPlacement="start" />
                <FormControlLabel control={<Switch />} label="Over 18" labelPlacement="start" />
                <FormControlLabel control={<Switch />} label="Søndag Dagspass Barn" labelPlacement="start" />
                <Typography variant="h3">Max Antall spillere: {maxNumberOfPlayers} </Typography>
            </FormGroup>
            <Paper elevation={2} sx={{ backgroundColor: 'rgpa(0,0,0,0.1)', padding: '1rem' }}>
                <AssigendPlayers />
            </Paper>

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
            <Typography variant="h2">Intreserte:</Typography>
            {poolPlayerInterests.map((poolPlayerInterest) => (
                <PlayerInterestInfo poolName={PoolName.saturdayEvening} playerInterest={poolPlayerInterest} />
            ))}
        </Paper>
    );
};

export default PlayerManagement;
