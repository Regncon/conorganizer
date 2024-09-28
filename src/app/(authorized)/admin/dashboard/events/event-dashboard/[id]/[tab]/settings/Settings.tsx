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
import { useEffect, useState, type FormEvent } from 'react';
import Slide from '@mui/material/Slide';
import Snackbar from '@mui/material/Snackbar';
import { Box, CircularProgress, Divider, Grid2, Stack, Switch } from '@mui/material';
import { onSnapshot, doc, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';
import { ConEvent } from '$lib/types';
import EventCardBig from '$app/(public)/components/components/EventCardBig';
import EventCardSmall from '$app/(public)/components/components/EventCardSmall';
import { updatePoolEvent } from './components/lib/actions';

type Props = {
    id: string;
};
const Settings = ({ id }: Props) => {
    const [user, setUser] = useState<User | null>();
    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);
    const [data, setData] = useState<ConEvent>();
    const eventDocRef = doc(db, 'events', id);

    const updateDatabase = async (data: Partial<ConEvent>) => {
        updateDoc(eventDocRef, data);
        console.log('updateDatabase', data);
        updatePoolEvent(id, data);
    };
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
    const handleSnackBar = (reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }
        setIsSnackBarOpen(false);
    };

    const handleOnChange = (e: FormEvent<HTMLFormElement>): void => {
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
            updatedBy: user.uid,
        };

        setIsSnackBarOpen(false);
        setSnackBarMessage('Endringar lagra!');
        setIsSnackBarOpen(true);

        updateDatabase(payload);
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
                        onChange={(evt) => handleOnChange(evt)}
                    >
                        <Grid2
                            size={{
                                xs: 12,
                                sm: 6,
                            }}
                        >
                            <Paper elevation={1} sx={{ padding: '1rem' }}>
                                <FormGroup sx={{ display: 'flex', gap: '1rem' }}>
                                    <FormLabel>Innstillinger</FormLabel>
                                    <FormControlLabel
                                        control={<Switch />}
                                        label="Publisert"
                                        name="published"
                                        checked={data.published}
                                        // onChange={() => setData({ ...data, published: !data.published })}
                                    />
                                    <FormControlLabel
                                        control={<Switch />}
                                        label="Kan arrangeres av andre frivillige"
                                        name="volunteersPossible"
                                        checked={data.volunteersPossible}
                                        // onChange={() => setData({ ...data, published: !data.published })}
                                    />
                                    <TextField
                                        sx={{ maxWidth: '15rem' }}
                                        type="number"
                                        name="participants"
                                        value={data.participants}
                                        // onChange={(e) => setData({ ...data, participants: parseInt(e.target.value) })}
                                        label="Maks antall deltakere"
                                        variant="outlined"
                                        required
                                    />
                                    <Stack component={FormControl} direction="row" spacing={1} alignItems="center">
                                        <Typography>Stor</Typography>
                                        <Switch
                                            name="isSmallCard"
                                            checked={data.isSmallCard}
                                            // onChange={() => setData({ ...data, isSmallCard: !data.isSmallCard })}
                                        />
                                        <Typography>Liten</Typography>
                                    </Stack>

                                    <Box sx={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
                                        <Box sx={{ opacity: data.isSmallCard ? '0.5' : 'unset', maxWidth: '297px' }}>
                                            <EventCardBig
                                                title={data.title}
                                                gameMaster={data.gameMaster}
                                                shortDescription={data.shortDescription}
                                                system={data.system}
                                                backgroundImage={data.smallImageURL}
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
                            <Paper sx={{ padding: '1rem', display: 'flex', gap: '1rem', flexDirection: 'column' }}>
                                <Typography variant="h2">Bilder</Typography>
                                <Typography>Filformat .webp</Typography>
                                <Typography>
                                    Tips: Ikke bruk spesialtegn i filnavn og unngå mellomrom og store bokstaver
                                </Typography>
                                <Typography component={'h4'}>Lite bilde bredde 430 og høyde 260</Typography>
                                <TextField
                                    type="text"
                                    name="smallImageURL"
                                    // onChange={(e) => setData({ ...data, smallImageURL: e.target.value })}
                                    value={data.smallImageURL}
                                    label="Lite bilde url"
                                    variant="outlined"
                                    fullWidth
                                />
                                <img src={data.smallImageURL ? data.smallImageURL : '/dice-small.webp'} alt="small" />
                                <Divider />

                                <Typography component={'h3'}>Stort blide bredde 1200 og høyde 212</Typography>
                                <TextField
                                    type="text"
                                    name="bigImageURL"
                                    // onChange={(e) => setData({ ...data, bigImageURL: e.target.value })}
                                    value={data.bigImageURL}
                                    label="Stort bilde url"
                                    variant="outlined"
                                    fullWidth
                                />
                                <Box
                                    component={'img'}
                                    maxWidth={430}
                                    src={data.bigImageURL ? data.bigImageURL : '/dice-big.webp'}
                                    alt="big"
                                />
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
                                        // onChange={(e) => setData({ ...data, email: e.target.value })}
                                        value={data.email}
                                        label="E-postadresse"
                                        variant="outlined"
                                        fullWidth
                                    />
                                    <TextField
                                        type="phone"
                                        name="phone"
                                        // onChange={(e) => setData({ ...data, phone: e.target.value })}
                                        value={data.phone}
                                        label="Telefonnummer"
                                        variant="outlined"
                                        fullWidth
                                    />
                                    <FormControlLabel
                                        control={<Checkbox name="moduleCompetition" checked={data.moduleCompetition} />}
                                        // onChange={(e) =>
                                        //     setData({ ...data, moduleCompetition: !data.moduleCompetition })
                                        // }
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
                                        // onChange={(e) => setData({ ...data, additionalComments: e.target.value })}
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
