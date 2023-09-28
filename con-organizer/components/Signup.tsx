'use client';

import { MouseEvent, useState } from 'react';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import { Alert, Box, Button, Card, IconButton, InputAdornment, Link, TextField } from '@mui/material';
import { FirebaseError } from 'firebase/app';
import { sendEmailVerification, signInWithEmailAndPassword } from 'firebase/auth';
import { auth } from '@/lib/firebase';

type Props = {
    setChoice: (choice: string) => void;
};

const Signup = ({ setChoice }: Props) => {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [passwordConfirmation, setPasswordConfirmation] = useState<string>('');
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const [success, setSuccess] = useState<string>('');
    const [error, setError] = useState<string>('');

    const handleClickShowPassword = () => setShowPassword((show) => !show);
    const handleMouseDownPassword = (event: MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };
    const signUp = async (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        if (password !== passwordConfirmation) {
            setError('Passordene matsjer isje, prøv igjen');
            return;
        }

        try {
            const signedInUser = await signInWithEmailAndPassword(auth, email, password);
            await sendEmailVerification(signedInUser.user);
            setSuccess('Suksess! Ett øyeblikk...');
            setError('');
            console.log(signedInUser);
            auth.signOut();
        } catch (e) {
            const err = e as FirebaseError;
            console.error(e);
            if (err.code === 'auth/invalid-email') {
                setError('Det ser ut som du skrev noe feil i epostadressen din, vennligst prøv igjen.');
            } else if (err.code === 'auth/too-many-requests') {
                setError(
                    'Vi har midlertidig suspendert kontoen din på grunn av for mange påloggingsforsøk. Dette er for at hackere og botter ikke skal kunne gjette seg til passordet ditt. Vennligst vent litt før du prøver igjen.'
                );
            } else {
                setError('Kunne ikke registrere deg fordi det skjedde en feil. Tekniske detaljer: ' + err.message);
            }
        }
        return;
    };
    return (
        <Box>
            <Card>
                <Box p={4} maxWidth={400} display="grid" justifyItems="center" gap={2}>
                    <form action={''}>
                        <TextField
                            label="e-post"
                            value={email}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                            fullWidth
                            sx={{ margin: '.3rem 0' }}
                        />
                        <TextField
                            sx={{ margin: '-1px 0 .3rem 0' }}
                            label="passord"
                            name="password"
                            value={password}
                            type={showPassword ? 'text' : 'password'}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            aria-label="veksle mellom synlig og skjult passord"
                                            onClick={handleClickShowPassword}
                                            onMouseDown={handleMouseDownPassword}
                                            edge="end"
                                        >
                                            {showPassword ? <VisibilityOff /> : <Visibility />}
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                            fullWidth
                        />
                        <TextField
                            sx={{ margin: '-1px 0 .3rem 0' }}
                            label="passord (igjen)"
                            name="password"
                            value={passwordConfirmation}
                            type={showPassword ? 'text' : 'password'}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            aria-label="veksle mellom synlig og skjult passord"
                                            onClick={handleClickShowPassword}
                                            onMouseDown={handleMouseDownPassword}
                                            edge="end"
                                        >
                                            {showPassword ? <VisibilityOff /> : <Visibility />}
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                setPasswordConfirmation(e.target.value)
                            }
                            fullWidth
                        />
                        <Button
                            variant="contained"
                            color="primary"
                            size="large"
                            type="submit"
                            fullWidth
                            onClick={(e) => signUp(e)}
                            sx={{ margin: '.3rem 0' }}
                        >
                            Lag konto
                        </Button>
                        <Button
                            variant="outlined"
                            size="large"
                            fullWidth
                            onClick={() => setChoice('')}
                            sx={{ margin: '-1px 0' }}
                        >
                            Avbryt
                        </Button>
                        {!success && !error
                            ? ''
                            : // <Alert severity="info">
                              //     OBS: Du kan ikke lage bruker uten &aring; ha kj&oslash;pt billett.&nbsp;
                              //     <Link href="https://www.regncon.no/kjop-billett-til-regncon-xxxi/">
                              //         Kjøp billett her!
                              //     </Link>
                              // </Alert>
                              null}
                        {!!success && <Alert severity="success">{success}</Alert>}
                        {!!error && <Alert severity="error">{error}</Alert>}
                    </form>
                </Box>
            </Card>
        </Box>
    );
};

export default Signup;
