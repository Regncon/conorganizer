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
import FormGroup from '@mui/material/FormGroup';
import Checkbox from '@mui/material/Checkbox';
import Button from '@mui/material/Button';
import Confetti from 'react-confetti';
import { useCallback, useEffect, useState, type ComponentProps, type FormEvent } from 'react';
import { doc, onSnapshot, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';

import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';
import Slide from '@mui/material/Slide';
import Skeleton from '@mui/material/Skeleton';
import Snackbar from '@mui/material/Snackbar';
import type { MyNewEvent } from '$lib/types';
type Props = {
    id: string;
};
const EventForm = ({ id }: Props) => {
    const [isExploding, setIsExploding] = useState(false);

    const [newEvent, setNewEvent] = useState<MyNewEvent>();
    const [user, setUser] = useState<User | null>();

    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);

    const newEventDocRef = doc(db, 'users', user?.uid ?? '_', 'my-events', id);

    const updateDatabase = async (newEvent: Partial<MyNewEvent>) => {
        await updateDoc(newEventDocRef, newEvent);
    };

    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (user) {
            unsubscribeSnapshot = onSnapshot(newEventDocRef, (snapshot) => {
                setNewEvent(snapshot.data() as MyNewEvent);
            });
        }

        const unsubscribeUser = onAuthStateChanged(firebaseAuth, (user) => {
            console.log(user, 'user');
            setUser(user);
        });

        return () => {
            unsubscribeSnapshot?.();
            unsubscribeUser();
        };
    }, [user]);

    const handleSnackBar = (event: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }

        setIsSnackBarOpen(false);
    };

    const handleBlur = useCallback(() => {
        setSnackBarMessage(snackBarMessageInitial);
        setIsSnackBarOpen(true);
    }, [isSnackBarOpen]);

    const handleSubmission = async () => {
        if (newEvent) {
            setIsSnackBarOpen(false);
            setIsExploding(!isExploding);
            await updateDatabase({ isSubmitted: !newEvent.isSubmitted });
            setSnackBarMessage(`Du har  ${newEvent.isSubmitted ? 'meldt av' : 'sendt inn'}  arrangementet`);
            setIsSnackBarOpen(true);
        }
    };
    const handleSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        e.preventDefault();
        if (!newEvent?.isSubmitted) {
            handleSubmission();
        }
    };

    const handleCancelSubmission: ComponentProps<'button'>['onClick'] = async (e) => {
        if (newEvent?.isSubmitted) {
            e.preventDefault();
            handleSubmission();
        }
    };

    const handleOnChange = (e: FormEvent<HTMLFormElement>) => {
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
        if (user?.email !== null) {
            let payload: Partial<MyNewEvent> = {
                [name]: value,
                updateAt: new Date(Date.now()).toString(),
                updatedBy: user?.email,
            };
            if (newEvent?.isSubmitted) {
                setIsSnackBarOpen(false);
                setSnackBarMessage('du må nå sende inn igjen skjemaet');
                payload = { ...payload, isSubmitted: false };
                setIsSnackBarOpen(true);
            }

            updateDatabase(payload);
        }
    };

    const skeletonHeight = 53;

    return newEvent ?
            <>
                <Grid2
                    sx={{ marginBlock: '1rem' }}
                    noValidate={newEvent.isSubmitted}
                    container
                    component="form"
                    spacing="2rem"
                    onBlur={handleBlur}
                    onChange={handleOnChange}
                    onSubmit={handleSubmit}
                >
                    {isExploding && (
                        <Confetti
                            onConfettiComplete={() => {
                                setIsExploding(!isExploding);
                            }}
                            numberOfPieces={2000}
                            recycle={false}
                            height={document.body.scrollHeight}
                            width={document.body.scrollWidth}
                        />
                    )}

                    <Grid2 xs={12}>
                        <Paper>
                            <FormLabel sx={{ padding: '1rem' }}>Du kan bruke opptil 33 teikn</FormLabel>
                            <TextField
                                name="title"
                                label="Tittel på spelmodul / arrangement"
                                value={newEvent.title}
                                variant="outlined"
                                required
                                fullWidth
                                inputProps={{
                                    title: 'Minst 1 teikn og maks teikn er 33',
                                }}
                                margin="dense"
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12}>
                        <Paper>
                            <FormLabel sx={{ padding: '1rem' }}>Du kan bruke opptil 50 teikn</FormLabel>
                            <TextField
                                name="subTitle"
                                value={newEvent.subTitle}
                                label="Kort oppsummering"
                                variant="outlined"
                                required
                                fullWidth
                                margin="dense"
                                inputProps={{
                                    title: 'Minst 1 teikn og maks 50 teikn',
                                }}
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12} sm={6} md={3}>
                        <Paper>
                            <TextField
                                type="email"
                                name="email"
                                value={newEvent.email}
                                label="E-postadresse"
                                variant="outlined"
                                required
                                fullWidth
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12} sm={6} md={3}>
                        <Paper>
                            <TextField
                                name="name"
                                value={newEvent.name}
                                label="Arrangørens namn (Ditt namn)"
                                variant="outlined"
                                required
                                fullWidth
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12} sm={6} md={3}>
                        <Paper>
                            <TextField
                                type="phone"
                                name="phone"
                                value={newEvent.phone}
                                label="Kva telefonnummer kan vi nå deg på?"
                                variant="outlined"
                                required
                                fullWidth
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12} sm={6} md={3}>
                        <Paper>
                            <TextField
                                name="system"
                                label="Spillsystem"
                                value={newEvent.system}
                                variant="outlined"
                                fullWidth
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl fullWidth>
                                <FormLabel>Skildring av modulen (tekst til programmet):</FormLabel>
                                <TextareaAutosize minRows={5} name="description" value={newEvent.description} />
                            </FormControl>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12} sm={4}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl fullWidth>
                                <FormLabel>Kva type spel er det?</FormLabel>
                                <RadioGroup
                                    value={newEvent.gameType}
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
                    <Grid2 xs={12} sm={4}>
                        <Paper sx={{ padding: '1rem', height: '100%' }}>
                            <TextField
                                type="number"
                                name="participants"
                                value={newEvent.participants}
                                label="Maks antall deltakere"
                                variant="outlined"
                                required
                                fullWidth
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12} sm={4}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormGroup>
                                <FormLabel>Kva for pulje kan du arrangere i?</FormLabel>
                                <FormControlLabel
                                    control={<Checkbox checked={newEvent.fridayEvening} />}
                                    name="fridayEvening"
                                    label="Fredag Kveld"
                                />
                                <FormControlLabel
                                    control={<Checkbox checked={newEvent.saturdayMorning} />}
                                    name="saturdayMorning"
                                    label="Lørdag Morgen"
                                />
                                <FormControlLabel
                                    control={<Checkbox checked={newEvent.saturdayEvening} />}
                                    name="saturdayEvening"
                                    label="Lørdag Kveld"
                                />
                                <FormControlLabel
                                    control={<Checkbox checked={newEvent.sundayMorning} />}
                                    name="sundayMorning"
                                    label="Søndag Morgen"
                                />
                            </FormGroup>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormGroup>
                                <FormLabel>Kryss av for det som gjeld</FormLabel>
                                <FormControlLabel
                                    control={<Checkbox name="moduleCompetition" checked={newEvent.moduleCompetition} />}
                                    label="Eg vil vere med på modulkonkurransen"
                                />
                                <Typography>
                                    husk å sende modulen til{' '}
                                    <a href="mailto:moduler@regncon.no ">moduler@regncon.no </a> innen første september!
                                </Typography>
                                <FormControlLabel
                                    control={<Checkbox name="childFriendly" checked={newEvent.childFriendly} />}
                                    label="Arrangementet passer for barn"
                                />
                                <FormControlLabel
                                    control={<Checkbox name="adultsOnly" checked={newEvent.adultsOnly} />}
                                    label="Arrangementet passer berre for vaksne (18+)"
                                />
                                <FormControlLabel
                                    control={<Checkbox name="beginnerFriendly" checked={newEvent.beginnerFriendly} />}
                                    label="Arrangementet er nybyrjarvenleg"
                                />
                                <FormControlLabel
                                    control={<Checkbox name="possiblyEnglish" checked={newEvent.possiblyEnglish} />}
                                    label="Arrangementet kan haldast på engelsk"
                                />
                                <FormControlLabel
                                    control={
                                        <Checkbox name="volunteersPossible" checked={newEvent.volunteersPossible} />
                                    }
                                    label="Andre kan halda arrangementet"
                                />
                                <FormControlLabel
                                    control={
                                        <Checkbox name="lessThanThreeHours" checked={newEvent.lessThanThreeHours} />
                                    }
                                    label="Eg trur arrangementet vil vare kortare enn 3 timer"
                                />
                                <FormControlLabel
                                    control={<Checkbox name="moreThanSixHours" checked={newEvent.moreThanSixHours} />}
                                    label="Eg trur arrangementet vil vare lenger enn 6 timer"
                                />
                            </FormGroup>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={12}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl fullWidth>
                                <FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
                                <TextareaAutosize
                                    minRows={3}
                                    name="additionalComments"
                                    value={newEvent.additionalComments}
                                />
                            </FormControl>
                        </Paper>
                    </Grid2>

                    <Grid2 xs={12}>
                        <Paper sx={{ padding: '1rem' }}>
                            <Typography>
                                Kladden vert lagra automatisk, men du må trykkje på knappen for å sende inn.
                            </Typography>
                            <Button type="submit" variant="contained" onClick={handleCancelSubmission}>
                                {newEvent.isSubmitted ? 'Meld av' : 'Send inn'}
                            </Button>
                        </Paper>
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
        :   <Grid2 container spacing="2rem">
                <Grid2 xs={12}>
                    <Skeleton variant="rounded" height={skeletonHeight} />
                </Grid2>
                <Grid2 xs={12}>
                    <Skeleton variant="rounded" height={skeletonHeight} />
                </Grid2>
                <Grid2 xs={12} sm={6} md={3}>
                    <Skeleton variant="rounded" height={skeletonHeight} />
                </Grid2>
                <Grid2 xs={12} sm={6} md={3}>
                    <Skeleton variant="rounded" height={skeletonHeight} />
                </Grid2>
                <Grid2 xs={12} sm={6} md={3}>
                    <Skeleton variant="rounded" height={skeletonHeight} />
                </Grid2>
                <Grid2 xs={12} sm={6} md={3}>
                    <Skeleton variant="rounded" height={skeletonHeight} />
                </Grid2>
                <Grid2 xs={12}>
                    <Skeleton variant="rounded" height={129} />
                </Grid2>
                <Grid2 xs={12} sm={4}>
                    <Skeleton variant="rounded" height={220} />
                </Grid2>
                <Grid2 xs={12} sm={4}>
                    <Skeleton variant="rounded" sx={{ height: { xs: skeletonHeight, sm: '220px' } }} />
                </Grid2>
                <Grid2 xs={12} sm={4}>
                    <Skeleton variant="rounded" height={220} />
                </Grid2>
                <Grid2 xs={12}>
                    <Skeleton variant="rounded" height={380} />
                </Grid2>
                <Grid2 xs={12}>
                    <Skeleton variant="rounded" height={90} />
                </Grid2>
                <Grid2 xs={12}>
                    <Skeleton variant="rounded" height={80} />
                </Grid2>
            </Grid2>;
};
export default EventForm;
