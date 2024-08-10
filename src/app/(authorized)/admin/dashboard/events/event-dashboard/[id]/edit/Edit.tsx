'use client';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import TextField from '@mui/material/TextField';
import FormControl from '@mui/material/FormControl';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import TextareaAutosize from '@mui/material/TextareaAutosize';
import Button from '@mui/material/Button';
import { useCallback, useEffect, useState, type FormEvent, type SyntheticEvent } from 'react';
import { ConEvent } from '$lib/types';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import Slide from '@mui/material/Slide';
import Snackbar, { type SnackbarCloseReason } from '@mui/material/Snackbar';
import Chip from '@mui/material/Chip';
import NavigateBeforeIcon from '@mui/icons-material/NavigateBefore';
import MainEvent from '$app/(public)/event/[id]/event';
import { Box, CircularProgress } from '@mui/material';
import { onSnapshot, doc, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';

type Props = {
    id: string;
};

const Edit = ({ id }: Props) => {
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
    const eventDocRef = doc(db, 'events', id);

    const [user, setUser] = useState<User | null>();
    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (user) {
            unsubscribeSnapshot = onSnapshot(eventDocRef, (snapshot) => {
                const newEventData = snapshot.data() as ConEvent;
                setData(newEventData);
                const newTags = [...tags].map((tag) => ({
                    ...tag,
                    selected: (newEventData[tag.name] as boolean) ?? false,
                }));
                setTags(newTags);
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

    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);
    const [tags, setTags] = useState<{ name: keyof ConEvent; label: string; selected: boolean }[]>([
        { name: 'childFriendly', label: 'Arrangementet passer for barn', selected: data?.childFriendly ?? false },
        {
            name: 'adultsOnly',
            label: 'Arrangementet passer berre for vaksne (18+)',
            selected: data?.adultsOnly ?? false,
        },
        {
            name: 'beginnerFriendly',
            label: 'Arrangementet er nybyrjarvenleg',
            selected: data?.beginnerFriendly ?? false,
        },
        {
            name: 'possiblyEnglish',
            label: 'Arrangementet kan haldast på engelsk',
            selected: data?.possiblyEnglish ?? false,
        },
        {
            name: 'volunteersPossible',
            label: 'Andre kan halda arrangementet',
            selected: data?.volunteersPossible ?? false,
        },
        {
            name: 'lessThanThreeHours',
            label: 'Eg trur arrangementet vil vare kortare enn 3 timer',
            selected: data?.lessThanThreeHours ?? false,
        },
        {
            name: 'moreThanSixHours',
            label: 'Eg trur arrangementet vil vare lenger enn 6 timer',
            selected: data?.moreThanSixHours ?? false,
        },
    ]);

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

    const updateDatabase = async (data: Partial<ConEvent>) => {
        await updateDoc(eventDocRef, data);
    };

    const handleSnackBar = (event: SyntheticEvent | globalThis.Event, reason?: SnackbarCloseReason) => {
        if (reason === 'clickaway') {
            return;
        }

        setIsSnackBarOpen(false);
    };

    const handleBlur = useCallback(() => {
        setSnackBarMessage(snackBarMessageInitial);
        setIsSnackBarOpen(true);
    }, [isSnackBarOpen]);

    // const handleSubmission = async () => {
    //     if (data) {
    //         setIsSnackBarOpen(false);
    //         setIsExploding(!isExploding);
    //         await updateDatabase({ isSubmitted: !data.isSubmitted });
    //         setSnackBarMessage(`Du har  ${data.isSubmitted ? 'meldt av' : 'sendt inn'}  arrangementet`);
    //         setIsSnackBarOpen(true);
    //     }
    // };
    // const handleSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
    //     e.preventDefault();
    //     if (!data?.isSubmitted) {
    //         handleSubmission();
    //     }
    // };
    //
    // const handleCancelSubmission: ComponentProps<'button'>['onClick'] = async (e) => {
    //     if (data?.isSubmitted) {
    //         e.preventDefault();
    //         handleSubmission();
    //     }
    // };
    const handleOnChange = useCallback(
        debounce((e: FormEvent<HTMLFormElement>): void => {
            const { value: inputValue, name: inputName, checked, type } = e.target as HTMLInputElement;

            let value: string | boolean = inputValue;
            let name: string = inputName;

            if (type === 'checkbox') {
                value = checked;
            }

            if (type === 'radio') {
                name = 'gameType';
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

    return (
        <>
            {!data ?
                <Typography variant="h1">
                    Loading...
                    <CircularProgress />
                </Typography>
                : <Box
                    component="form"
                    onChange={(evt) =>
                        handleOnChange(evt).catch((err) => {
                            if (err !== 'Aborted by debounce') {
                                throw err;
                            }
                        })
                    }
                >
                    <MainEvent id={id} editable={true} />
                    <Snackbar
                        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                        open={isSnackBarOpen}
                        onClose={handleSnackBar}
                        TransitionComponent={Slide}
                        message={snackBarMessage}
                        autoHideDuration={1200}
                    />
                </Box>
            }
        </>
    );
    /*
    return (
        <>
            <Grid2
                sx={{ marginBlock: '1rem' }}
                noValidate={data.isSubmitted}
                container
                component="form"
                spacing="2rem"
                onBlur={handleBlur}
                onChange={handleOnChange}
                onSubmit={handleSubmit}
            >
                <Grid2 xs={5}>
                    <MainEvent eventData={event} editable={true} />
                </Grid2>
                <Grid2>
                    <Grid2 xs={5}>
                        <Paper>
                            <FormLabel sx={{ padding: '1rem' }}>Du kan bruke opptil 33 teikn</FormLabel>
                            <TextField
                                name="title"
                                label="Tittel på spelmodul / arrangement"
                                value={data.title}
                                variant="outlined"
                                required
                                inputProps={{
                                    title: 'Minst 1 teikn og maks teikn er 33',
                                }}
                                margin="dense"
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper>
                            <FormLabel sx={{ padding: '1rem' }}>Du kan bruke opptil 50 teikn</FormLabel>
                            <TextField
                                name="subTitle"
                                value={data.subTitle}
                                label="Kort oppsummering"
                                variant="outlined"
                                required
                                margin="dense"
                                inputProps={{
                                    title: 'Minst 1 teikn og maks 50 teikn',
                                }}
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6} sm={6} md={3}>
                        <Paper>
                            <TextField
                                name="name"
                                value={data.name}
                                label="Arrangørens namn (Ditt namn)"
                                variant="outlined"
                                required
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6} sm={6} md={3}>
                        <Paper>
                            <TextField name="system" label="Spillsystem" value={data.system} variant="outlined" />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl>
                                <FormLabel>Skildring av modulen (tekst til programmet):</FormLabel>
                                <TextareaAutosize minRows={5} name="description" value={data.description} />
                            </FormControl>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6} sm={4}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl>
                                <FormLabel>Kva type spel er det?</FormLabel>
                                <RadioGroup
                                    value={data.gameType}
                                    aria-labelledby="demo-controlled-radio-buttons-group"
                                    name="controlled-radio-buttons-group"
                                >
                                    <FormControlLabel
                                        value="rolePlaying"
                                        control={<Radio name="rolePlaying" />}
                                        label="rollespel"
                                    />
                                    <FormControlLabel
                                        value="boardGame"
                                        control={<Radio name="boardGame" />}
                                        label="Brettspel"
                                    />
                                    <FormControlLabel
                                        value="cardGame"
                                        control={<Radio name="cardGame" />}
                                        label="Kortspel"
                                    />
                                    <FormControlLabel value="other" control={<Radio name="other" />} label="Annet" />
                                </RadioGroup>
                            </FormControl>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <Typography sx={{ marginBlockStart: '1rem' }}>
                                Trykk på brikka som passar til spelet ditt:
                            </Typography>
                            <Box sx={{ maxWidth: '320px' }}>
                                {tags.map((tag) => (
                                    <Chip
                                        sx={{ marginBlock: '0.4rem', marginInlineEnd: '0.4rem' }}
                                        label={tag.label}
                                        key={tag.name}
                                        color={tag.selected ? 'primary' : 'secondary'}
                                        icon={<NavigateBeforeIcon />}
                                        variant={tag.selected ? 'filled' : 'outlined'}
                                        onClick={async () => {
                                            setTags((prev) =>
                                                prev.map((t) =>
                                                    t.name === tag.name ? { ...t, selected: !t.selected } : t
                                                )
                                            );
                                            await updateDatabase({ [tag.name]: !tag.selected });
                                            setIsSnackBarOpen(true);
                                        }}
                                    />
                                ))}
                            </Box>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl>
                                <FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
                                <TextareaAutosize
                                    minRows={3}
                                    name="additionalComments"
                                    value={data.additionalComments}
                                />
                            </FormControl>
                        </Paper>
                    </Grid2>

                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <Typography>
                                Kladden vert lagra automatisk, men du må trykkje på knappen for å sende inn.
                            </Typography>
                            <Button type="submit" variant="contained" onClick={handleCancelSubmission}>
                                {data.isSubmitted ? 'Meld av' : 'Send inn'}
                            </Button>
                        </Paper>
                    </Grid2>
                </Grid2>
            </Grid2>
            <Snackbar
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                open={isSnackBarOpen}
                onClose={handleSnackBar}
                TransitionComponent={Slide}
                message={snackBarMessage}
                autoHideDuration={1200}
            />
        </>
    );
        */
};
export default Edit;
