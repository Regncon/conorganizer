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
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import Slide from '@mui/material/Slide';
import Snackbar from '@mui/material/Snackbar';
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    CircularProgress,
    IconButton,
    Radio,
    RadioGroup,
    Stack,
    Switch,
} from '@mui/material';
import WarningIcon from '@mui/icons-material/Warning';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ArrowDropUpIcon from '@mui/icons-material/ArrowDropUp';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { onSnapshot, doc, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';
import { ConEvent, Pulje } from '$lib/types';

type Props = {
    id: string;
};
const Settings = ({ id }: Props) => {
    /**
     * Debounces a function, creating a new function that does the same as the original, but will not actually run before
     * a specified amount of time has passed since it was last called.
     *
     * @param fn The function to debounce
     * @param delay Number of milliseconds to wait since the last call to the function to actually run it
     *
     * @returns A function that does the same as `fn`, but won't actually run before `delay` milliseconds has passed since
     * its last invocation. Its return value will be wrapped in a promise
     */
    const debounce = <P extends unknown[], R>(
        fn: (...args: P) => R | Promise<R>,
        delay: Parameters<typeof setTimeout>[1]
    ): ((...args: P) => Promise<R>) => {
        let timer: ReturnType<typeof setTimeout> | null = null;

        type Reject = Parameters<ConstructorParameters<typeof Promise<R>>[0]>[1];
        let prevReject: Reject = () => { };

        return (...args: P): Promise<R> =>
            new Promise((resolve, reject) => {
                if (timer !== null) {
                    clearTimeout(timer);
                    prevReject('Aborted by debounce');
                }

                prevReject = reject;

                timer = setTimeout(async () => {
                    timer = null;

                    try {
                        const result = await fn(...args);
                        resolve(result);
                    } catch (err) {
                        reject(err);
                    }
                }, delay);
            });
    };
    const [user, setUser] = useState<User | null>();
    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);

    const initialState: ConEvent = {
        gameMaster: '',
        id: '',
        shortDescription: '',
        description: '',
        system: '',
        title: '',
        email: '',
        name: '',
        phone: '',
        gameType: '',
        participants: 0,
        unwantedFridayEvening: false,
        unwantedSaturdayMorning: false,
        unwantedSaturdayEvening: false,
        unwantedSundayMorning: false,
        moduleCompetition: false,
        childFriendly: false,
        possiblyEnglish: false,
        adultsOnly: false,
        volunteersPossible: false,
        lessThanThreeHours: false,
        moreThanSixHours: false,
        beginnerFriendly: false,
        additionalComments: '',
        createdAt: '',
        createdBy: '',
        updateAt: '',
        updatedBy: '',
        subTitle: '',
        published: false,
        puljeFridayEvening: false,
        puljeSaturdayMorning: false,
        puljeSaturdayEvening: false,
        puljeSundayMorning: false,
    };
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

    const handleBlur = useCallback(() => {
        setSnackBarMessage(snackBarMessageInitial);
        setIsSnackBarOpen(true);
    }, [isSnackBarOpen]);

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

    let eventOrder = [
        { id: '0', name: 'Opprop fredag kveld Kl: 17:30', order: 0 },
        { id: '1', name: 'Orker på tur', order: 1 },
        { id: '2', name: 'Cathulu er forelsket', order: 2, thisEvent: true },
    ];

    return (
        <>
            {!data ?
                <Typography variant="h1">
                    Loading...
                    <CircularProgress />
                </Typography>
                : <>
                    <Grid2
                        sx={{ padding: '1rem' }}
                        container
                        component="form"
                        spacing="2rem"
                        onChange={(evt) =>
                            handleOnChange(evt).catch((err) => {
                                if (err !== 'Aborted by debounce') {
                                    throw err;
                                }
                            })
                        }
                    >
                        <Grid2 xs={12} sm={6}>
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
                                    <Paper elevation={3}>
                                        <Accordion>
                                            <AccordionSummary
                                                expandIcon={<ExpandMoreIcon />}
                                                aria-controls="panel1-content"
                                                id="panel1-header"
                                            >
                                                Sortering rekkefølge
                                            </AccordionSummary>
                                            <AccordionDetails>
                                                <Typography variant={'h3'} sx={{ textAlign: 'center' }}>
                                                    Pulje: Fredag kveld
                                                </Typography>
                                                {eventOrder.map((event) => (
                                                    <Paper
                                                        key={event.id}
                                                        elevation={4}
                                                        sx={{
                                                            padding: '1rem',
                                                            marginBottom: '1rem',
                                                            display: 'flex',
                                                            justifyContent: 'space-between',
                                                            backgroundColor: event.thisEvent ? 'primary.light' : '',
                                                        }}
                                                    >
                                                        <Typography component={'span'}> {event.name}</Typography>
                                                        <Box
                                                            sx={{
                                                                display: 'inline-block',
                                                            }}
                                                        >
                                                            <IconButton>
                                                                <ArrowDropUpIcon />
                                                            </IconButton>
                                                            <IconButton>
                                                                <ArrowDropDownIcon />
                                                            </IconButton>
                                                        </Box>
                                                    </Paper>
                                                ))}
                                            </AccordionDetails>
                                        </Accordion>
                                    </Paper>
                                </FormGroup>
                            </Paper>
                        </Grid2>
                        <Grid2 xs={12} sm={6}>
                            <Paper sx={{ padding: '1rem' }}>
                                <FormGroup>
                                    <FormLabel>Pulje</FormLabel>
                                    <FormControlLabel
                                        control={
                                            <Checkbox name="puljeFridayEvening" checked={data.puljeFridayEvening} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeFridayEvening: !data.puljeFridayEvening })
                                        }
                                        label={unwantedTimeSlotWarning('Fredag Kveld', data.unwantedFridayEvening)}
                                    />
                                    <FormControlLabel
                                        control={
                                            <Checkbox name="puljeSaturdayMorning" checked={data.puljeSaturdayMorning} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeSaturdayMorning: !data.puljeSaturdayMorning })
                                        }
                                        label={unwantedTimeSlotWarning('Lørdag Morgen', data.unwantedSaturdayMorning)}
                                    />
                                    <FormControlLabel
                                        control={
                                            <Checkbox name="puljeSaturdayEvening" checked={data.puljeSaturdayEvening} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeSaturdayEvening: !data.puljeSaturdayEvening })
                                        }
                                        label={unwantedTimeSlotWarning('Lørdag Kveld', data.unwantedSaturdayEvening)}
                                    />
                                    <FormControlLabel
                                        control={
                                            <Checkbox name="puljeSundayMorning" checked={data.puljeSundayMorning} />
                                        }
                                        onChange={() =>
                                            setData({ ...data, puljeSundayMorning: !data.puljeSundayMorning })
                                        }
                                        label={unwantedTimeSlotWarning('Søndag Morgen', data.unwantedSundayMorning)}
                                    />
                                </FormGroup>
                            </Paper>
                        </Grid2>
                        <Grid2 xs={12} sm={6}>
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
                        <Grid2 xs={12}>
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
