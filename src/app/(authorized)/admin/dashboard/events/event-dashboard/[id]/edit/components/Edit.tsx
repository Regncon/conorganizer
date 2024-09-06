'use client';
import Typography from '@mui/material/Typography';
import { useCallback, useEffect, useState, type FormEvent, type SyntheticEvent } from 'react';
import { ConEvent } from '$lib/types';
import Slide from '@mui/material/Slide';
import Snackbar, { type SnackbarCloseReason } from '@mui/material/Snackbar';
import { Box, CircularProgress } from '@mui/material';
import { onSnapshot, doc, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';
import DescriptionDialog from './DescriptionDialog';
import debounce from '$lib/debounce';
import MainEvent from '$app/(public)/event/[id]/components/MainEvent';

type Props = {
    id: string;
};

const Edit = ({ id }: Props) => {
    const initialState: ConEvent = {} as ConEvent;
    const [data, setData] = useState<ConEvent>();
    const eventDocRef = doc(db, 'events', id);

    const [user, setUser] = useState<User | null>();
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

    const [openDescriptionDialog, setOpenDescriptionDialog] = useState<boolean>(false);
    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);

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
            saveToDb(name, value);
        }, 1500),
        [user]
    );

    const saveToDb = (name: string, value: string | boolean) => {
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
    };

    return (
        <>
            {!data ?
                <Typography variant="h1">
                    Loading...
                    <CircularProgress />
                </Typography>
            :   <>
                    <Box
                        component="form"
                        onChange={(evt) =>
                            handleOnChange(evt).catch((err) => {
                                if (err !== 'Aborted by debounce') {
                                    throw err;
                                }
                            })
                        }
                    >
                        <MainEvent id={id} editable={true} editDescription={(edit) => setOpenDescriptionDialog(edit)} />
                        <Snackbar
                            anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                            open={isSnackBarOpen}
                            onClose={handleSnackBar}
                            TransitionComponent={Slide}
                            message={snackBarMessage}
                            autoHideDuration={1200}
                        />
                    </Box>
                    <DescriptionDialog
                        data={data}
                        handleSave={() => {
                            saveToDb('description', data.description);
                            setOpenDescriptionDialog(false);
                        }}
                        open={openDescriptionDialog}
                        close={() => setOpenDescriptionDialog(false)}
                    />
                </>
            }
        </>
    );
};
export default Edit;
