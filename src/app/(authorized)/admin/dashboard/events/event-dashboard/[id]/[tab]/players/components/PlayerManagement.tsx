import {
    Box,
    Button,
    FormControl,
    FormControlLabel,
    FormGroup,
    Input,
    InputAdornment,
    InputLabel,
    Link,
    Paper,
    Switch,
    Typography,
} from '@mui/material';
import { InterestLevel, PoolName, RoomName } from '$lib/enums';
import SearchIcon from '@mui/icons-material/Search';
import { generatePoolPlayerInterestById } from './lib/actions';
import PlayerInterestInfo from './components/PlayerInterestInfo';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';
import AssignedPlayers from './components/AssignedPlayers';
import { PlayerInterest } from '$lib/types';
import ManuallyAssignPlayer from './components/manuallyAssignPlayer';
import NextLink from 'next/link';

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
    // show isGameMaster at the top of the list
    const assignedPlayers = poolPlayerInterests
        .filter((player) => player.isAssigned)
        .sort((a, b) => (a.isGameMaster ? -1 : 1));
    console.log('assignedPlayers', assignedPlayers);

    // console.log('poolPlayerInterests', poolPlayerInterests);

    const dummyPlayerInrerestData: PlayerInterest = {
        interestLevel: InterestLevel.VeryInterested,
        poolEventId: '',
        participantId: '',
        firstName: 'Kari',
        lastName: 'Nordmann',
        isOver18: true,
        ticketCategoryID: 0,
        ticketCategory: 'Festivalpass Ungdom/student (13-25/30år) Early-bird',
        playerInPools: [
            {
                id: undefined,
                participantId: '',
                firstName: '',
                lastName: '',
                interestLevel: InterestLevel.VeryInterested,
                poolEventId: '',
                poolEventTitle: 'Et gøyalt spill',
                poolName: PoolName.fridayEvening,
                isFirstChoice: true,
                isGameMaster: false,
                createdAt: '',
                createdBy: '',
                updateAt: '',
                updatedBy: '',
                roomId: '',
                roomName: RoomName.NotSet,
                isPublished: false,
                isAssigned: false,
            },
            {
                id: undefined,
                participantId: '',
                firstName: '',
                lastName: '',
                interestLevel: InterestLevel.SomewhatInterested,
                poolEventId: '',
                poolEventTitle: '"Hva f**n skjedde!?!??" -- ghost/echo av John Harper',
                poolName: PoolName.saturdayMorning,
                isFirstChoice: false,
                isGameMaster: true,
                createdAt: '',
                createdBy: '',
                updateAt: '',
                updatedBy: '',
                roomId: '',
                roomName: RoomName.NotSet,
                isPublished: false,
                isAssigned: false,
            },
            {
                id: undefined,
                participantId: '',
                firstName: '',
                lastName: '',
                interestLevel: InterestLevel.Interested,
                poolEventId: '',
                poolEventTitle: 'Random dnd event',
                poolName: PoolName.saturdayMorning,
                isFirstChoice: false,
                isGameMaster: true,
                createdAt: '',
                createdBy: '',
                updateAt: '',
                updatedBy: '',
                roomId: '',
                roomName: RoomName.NotSet,
                isPublished: false,
                isAssigned: false,
            },
        ],
        isGameMaster: false,
        isAssigned: false,
        currentPoolPlayerId: undefined,
        isAlredyPlayerInPool: false,
    };
    const dummyPlayersInrerestData: PlayerInterest[] = [
        dummyPlayerInrerestData,
        dummyPlayerInrerestData,
        dummyPlayerInrerestData,
    ];
    return (
        <Paper sx={{ padding: '1rem' }}>
            <Paper
                component={FormGroup}
                sx={{
                    display: 'flex',
                    flexDirection: 'row',
                    gap: '1rem',
                    placeItems: 'center',
                    position: 'sticky',
                    top: 'var(--app-bar-height-desktop)',
                    zIndex: 1,
                }}
            >
                <Typography variant="h1">{translatedDays.get(poolName)}</Typography>
                <FormControlLabel control={<Switch />} label="Under 18" labelPlacement="start" />
                <FormControlLabel control={<Switch />} label="Over 18" labelPlacement="start" />
                <FormControlLabel control={<Switch />} label="Søndag Dagspass Barn" labelPlacement="start" />
                <Typography variant="h3">
                    Antall tildelte spillere:{' '}
                    {assignedPlayers.filter((assignedPlayer) => !assignedPlayer.isGameMaster).length}
                </Typography>
                <Typography variant="h3">Max Antall spillere: {maxNumberOfPlayers} </Typography>
                <Box sx={{ display: 'flex', width: '100%', gap: '1rem', placeItems: 'center', marginBlock: '1rem' }}>
                    <Button
                        sx={{ maxWidth: 'fit-content', maxHeight: 'fit-content' }}
                        variant="contained"
                        color="primary"
                    >
                        Algotitme
                    </Button>
                    <Button
                        sx={{ maxWidth: 'fit-content', maxHeight: 'fit-content' }}
                        variant="contained"
                        color="primary"
                    >
                        Algotitme2
                    </Button>
                    <Button
                        sx={{ maxWidth: 'fit-content', maxHeight: 'fit-content' }}
                        variant="contained"
                        color="primary"
                    >
                        Algotitme3
                    </Button>
                    <Box sx={{ placeSelf: 'end' }}>
                        <Link component={NextLink} href={'#assigned-players'} color="secondary">
                            tilbake til toppen
                        </Link>
                    </Box>
                </Box>
            </Paper>
            <Paper elevation={2} sx={{ backgroundColor: 'rgpa(0,0,0,0.1)', padding: '1rem' }}>
                <AssignedPlayers poolName={poolName} assignedPlayers={assignedPlayers} />
            </Paper>

            <ManuallyAssignPlayer poolEventId={id} />
            <Typography variant="h2">Intreserte:</Typography>
            {poolPlayerInterests
                .sort((a, b) =>
                    a.isGameMaster ? -1
                    : a.interestLevel > b.interestLevel ? -1
                    : 1
                )
                .map((poolPlayerInterest, index) => (
                    <PlayerInterestInfo
                        key={index}
                        poolName={poolName}
                        playerInterest={poolPlayerInterest}
                        hasPlayerInPool={poolPlayerInterest.isAlredyPlayerInPool}
                    />
                ))}
        </Paper>
    );
};

export default PlayerManagement;
