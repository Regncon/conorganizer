'use client';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import TextField from '@mui/material/TextField';
import FormControl from '@mui/material/FormControl';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import TextareaAutosize from '@mui/material/TextareaAutosize';
import FormGroup from '@mui/material/FormGroup';
import Checkbox from '@mui/material/Checkbox';
import { useCallback, useEffect, useState, type FormEvent } from 'react';
import Slide from '@mui/material/Slide';
import Snackbar from '@mui/material/Snackbar';
import { Box, CircularProgress, Grid2, Stack, Switch } from '@mui/material';
import { onSnapshot, doc, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';
import { ConEvent } from '$lib/types';
import Ordering from './ordering';
import debounce from '$lib/debounce';
import WarningIcon from '@mui/icons-material/Warning';
import EventCardBig from '$app/(public)/components/EventCardBig';
import EventCardSmall from '$app/(public)/components/EventCardSmall';

type Props = {
    id: string;
    allEvents: ConEvent[];
};
const Settings = ({ id, allEvents }: Props) => {
    const [user, setUser] = useState<User | null>();
    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);
    const initialState = {} as ConEvent;
    const [data, setData] = useState<ConEvent>();
    // console.log('data', data);
    const eventDocRef = doc(db, 'events', id);

    const updateDatabase = async (data: Partial<ConEvent>) => {
        await updateDoc(eventDocRef, data);
    };
    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (id !== undefined) {
            unsubscribeSnapshot = onSnapshot(doc(db, 'events', id), (snapshot) => {
                setData((snapshot.data() as ConEvent | undefined) ?? initialState);
            });
        }
        return () => {
            unsubscribeSnapshot?.();
        };
    }, [id]);

    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (user) {
            unsubscribeSnapshot = onSnapshot(eventDocRef, (snapshot) => {
                const newEventData = snapshot.data() as ConEvent;
                setData(newEventData);
            });
        }

        const unsubscribeUser = onAuthStateChanged(firebaseAuth, (user) => {
            setUser(user);
        });

        return () => {
            unsubscribeSnapshot?.();
            unsubscribeUser();
        };
    }, [user]);

    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (id !== undefined) {
            unsubscribeSnapshot = onSnapshot(doc(db, 'events', id), (snapshot) => {
                setData((snapshot.data() as ConEvent | undefined) ?? initialState);
            });
        }
        return () => {
            unsubscribeSnapshot?.();
        };
    }, [id]);
    const handleSnackBar = (reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }
        setIsSnackBarOpen(false);
    };

    const handleOnChange = useCallback(
        debounce((e: FormEvent<HTMLFormElement>): void => {
            const { value: inputValue, name: inputName, checked, type } = e.target as HTMLInputElement;

            let value: string | boolean = inputValue;
            let name: string = inputName;

            if (type === 'checkbox') {
                value = checked;
            }

            if (type === 'radio') {
                name = 'pulje';
                value = inputName;
            }
            if (!user || !user.email) {
                console.error('user?.email is null');
                return;
            }

            let payload: Partial<ConEvent> = {
                [name]: value,
                updateAt: new Date(Date.now()).toString(),
                updatedBy: user.email,
            };
            setIsSnackBarOpen(false);
            setSnackBarMessage('Endringar lagra!');
            setIsSnackBarOpen(true);

            updateDatabase(payload);
        }, 1500),
        [user]
    );

    const unwantedTimeSlotWarning = (slot: string, unwanted: boolean) => {
        return (
            <Stack direction="row">
                <Typography component={'span'} sx={{ paddingRight: '1rem' }}>
                    {slot}
                </Typography>
                {unwanted && (
                    <Box sx={{ display: 'inherit', color: 'warning.main' }}>
                        <WarningIcon />
                        <Typography component={'i'}>Gm ønsker ikke denne</Typography>
                    </Box>
                )}
            </Stack>
        );
    };

    return (
        <>
            {!data ?
                <Typography variant="h1">
                    Loading...
                    <CircularProgress />
                </Typography>
            :   <>
                    <Grid2
                        sx={{
                            paddingTop: '1rem',
                        }}
                        container
                        component="form"
                        rowSpacing={{ xs: 1, md: 2 }}
                        columnSpacing={{ xs: 0, sm: 1, md: 2 }}
                        onChange={(evt) =>
                            handleOnChange(evt).catch((err) => {
                                if (err !== 'Aborted by debounce') {
                                    throw err;
                                }
                            })
                        }
                    >
                        <Grid2
                            size={{
                                xs: 12,
                                sm: 6,
                            }}
                        >
                            <Paper elevation={1} sx={{ padding: '1rem' }}>
                                <FormGroup sx={{ display: 'flex', gap: '1rem' }}>
                                    <FormLabel>Instillinger</FormLabel>
                                    <FormControlLabel
                                        control={<Switch />}
                                        label="Publisert"
                                        name="published"
                                        checked={data.published}
                                        onChange={() => setData({ ...data, published: !data.published })}
                                    />
                                    <TextField
                                        sx={{ maxWidth: '15rem' }}
                                        type="number"
                                        name="participants"
                                        value={data.participants}
                                        onChange={(e) => setData({ ...data, participants: parseInt(e.target.value) })}
                                        label="Maks antall deltakere"
                                        variant="outlined"
                                        required
                                    />
                                    <Stack component={FormControl} direction="row" spacing={1} alignItems="center">
                                        <Typography>Stor</Typography>
                                        <Switch
                                            name="isSmallCard"
                                            checked={data.isSmallCard}
                                            onChange={() => setData({ ...data, isSmallCard: !data.isSmallCard })}
                                        />
                                        <Typography>Liten</Typography>
                                    </Stack>

                                    <Box sx={{ display: 'flex', gap: '1rem' }}>
                                        <Box sx={{ opacity: data.isSmallCard ? '0.5' : 'unset' }}>
                                            <EventCardBig
                                                title={data.title}
                                                gameMaster={data.gameMaster}
                                                shortDescription={data.shortDescription}
                                                system={data.system}
                                            />
                                        </Box>
                                        <Box sx={{ opacity: !data.isSmallCard ? '0.5' : 'unset' }}>
                                            <EventCardSmall
                                                title={data.title}
                                                gameMaster={data.gameMaster}
                                                system={data.system}
                                            />
                                        </Box>
                                    </Box>
                                </FormGroup>
                            </Paper>
                        </Grid2>
                        <Grid2
                            size={{
                                xs: 12,
                                sm: 6,
                            }}
                        >
                            <Paper sx={{ padding: '1rem' }}>
                                <FormGroup>
                                    <FormLabel>Pulje</FormLabel>
                                    <FormControlLabel
                                        disabled
                                        control={
                                            <Checkbox name="puljeFridayEvening" checked={data.puljeFridayEvening} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeFridayEvening: !data.puljeFridayEvening })
                                        }
                                        label={unwantedTimeSlotWarning('Fredag Kveld', data.unwantedFridayEvening)}
                                    />
                                    <Ordering
                                        id={id}
                                        allEvents={allEvents}
                                        pulje="Fredag kveld"
                                        disabled={!data.puljeFridayEvening}
                                    />
                                    <FormControlLabel
                                        disabled
                                        control={
                                            <Checkbox name="puljeSaturdayMorning" checked={data.puljeSaturdayMorning} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeSaturdayMorning: !data.puljeSaturdayMorning })
                                        }
                                        label={unwantedTimeSlotWarning('Lørdag Morgen', data.unwantedSaturdayMorning)}
                                    />
                                    <Ordering
                                        id={id}
                                        allEvents={allEvents}
                                        pulje="Lørdag morgen"
                                        disabled={!data.puljeSaturdayMorning}
                                    />
                                    <FormControlLabel
                                        disabled
                                        control={
                                            <Checkbox name="puljeSaturdayEvening" checked={data.puljeSaturdayEvening} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeSaturdayEvening: !data.puljeSaturdayEvening })
                                        }
                                        label={unwantedTimeSlotWarning('Lørdag Kveld', data.unwantedSaturdayEvening)}
                                    />
                                    <Ordering
                                        id={id}
                                        allEvents={allEvents}
                                        pulje="Lørdag kveld"
                                        disabled={!data.puljeSaturdayEvening}
                                    />
                                    <FormControlLabel
                                        disabled
                                        control={
                                            <Checkbox name="puljeSundayMorning" checked={data.puljeSundayMorning} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeSundayMorning: !data.puljeSundayMorning })
                                        }
                                        label={unwantedTimeSlotWarning('Søndag Morgen', data.unwantedSundayMorning)}
                                    />
                                    <Ordering
                                        id={id}
                                        allEvents={allEvents}
                                        pulje="Søndag morgen"
                                        disabled={!data.puljeSundayMorning}
                                    />
                                </FormGroup>
                            </Paper>
                        </Grid2>
                        <Grid2
                            size={{
                                xs: 12,
                                sm: 6,
                            }}
                        >
                            <Paper sx={{ padding: '1rem' }}>
                                <FormGroup sx={{ display: 'flex', gap: '1rem' }}>
                                    <FormLabel>Kontaktinfo</FormLabel>

                                    <TextField
                                        type="email"
                                        name="email"
                                        onChange={(e) => setData({ ...data, email: e.target.value })}
                                        value={data.email}
                                        label="E-postadresse"
                                        variant="outlined"
                                        fullWidth
                                    />
                                    <TextField
                                        type="phone"
                                        name="phone"
                                        onChange={(e) => setData({ ...data, phone: e.target.value })}
                                        value={data.phone}
                                        label="Telefonnummer"
                                        variant="outlined"
                                        fullWidth
                                    />
                                    <FormControlLabel
                                        control={<Checkbox name="moduleCompetition" checked={data.moduleCompetition} />}
                                        onChange={(e) =>
                                            setData({ ...data, moduleCompetition: !data.moduleCompetition })
                                        }
                                        label="Modulen er påmeldt konkurransen"
                                    />
                                </FormGroup>
                            </Paper>
                        </Grid2>
                        <Grid2 size={12}>
                            <Paper sx={{ padding: '1rem' }}>
                                <FormControl fullWidth>
                                    <FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
                                    <TextareaAutosize
                                        minRows={6}
                                        name="additionalComments"
                                        value={data.additionalComments}
                                        onChange={(e) => setData({ ...data, additionalComments: e.target.value })}
                                    />
                                </FormControl>
                            </Paper>
                        </Grid2>
                    </Grid2>
                    <Snackbar
                        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                        open={isSnackBarOpen}
                        onClose={(e, r?) => handleSnackBar(r)}
                        TransitionComponent={Slide}
                        message={snackBarMessage}
                        autoHideDuration={1200}
                    />
                </>
            }
        </>
    );
};
export default Settings;
