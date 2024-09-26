'use client';
import React from 'react';
import Button from '@mui/material/Button';
import SvgIcon from '@mui/material/SvgIcon';
import { firebaseAuth } from '$lib/firebase/firebase';
import { signInWithPopup, GoogleAuthProvider } from 'firebase/auth';
import { useRouter } from 'next/navigation';
import type { Route } from 'next';
import { GetMyParticipants } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';

const GoogleIcon = () => (
    <SvgIcon viewBox="0 0 48 48">
        <path
            fill="#EA4335"
            d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"
        />
        <path
            fill="#4285F4"
            d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"
        />
        <path
            fill="#FBBC05"
            d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z"
        />
        <path
            fill="#34A853"
            d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"
        />
        <path fill="none" d="M0 0h48v48H0z" />
    </SvgIcon>
);

type Props = {
    redirectTo?: Route;
    disabled?: boolean;
};

const GoogleSignInButton = ({ redirectTo = '/', disabled = false }: Props) => {
    const router = useRouter();
    const handleClick = async () => {
        const provider = new GoogleAuthProvider();
        signInWithPopup(firebaseAuth, provider)
            .then(async (result) => {
                router.prefetch(redirectTo);
                // This gives you a Google Access Token. You can use it to access the Google API.
                const credential = GoogleAuthProvider.credentialFromResult(result);
                const token = credential?.accessToken;
                console.log('token', token);
                const idToken = credential?.idToken;
                console.log('idToken', credential?.idToken);
                //await setSessionCookie(idToken?? '');

                // The signed-in user info.
                const user = result.user;
                // IdP data available using getAdditionalUserInfo(result)
                // ...
                //
                GetMyParticipants().then((myParticipants) => {
                    localStorage.setItem('myParticipants', JSON.stringify(myParticipants));
                    router.replace(redirectTo);
                });
            })
            .catch((error) => {
                // Handle Errors here.
                const errorCode = error.code;
                console.error('errorCode', errorCode);
                const errorMessage = error.message;
                console.error('errorMessage', errorMessage);
                // The email of the user's account used.
                const email = error.customData?.email;
                console.error('email', email);
                // The AuthCredential type that was used.
                const credential = GoogleAuthProvider.credentialFromError(error);
                console.error('credential', credential);
                // ...
            });
    };
    return (
        <Button
            onClick={handleClick}
            variant="outlined"
            startIcon={<GoogleIcon />}
            sx={{
                placeSelf: 'center',
                backgroundColor: '#131314',
                borderColor: '#8e918f',
                color: '#e3e3e3',
                borderRadius: '20px',
                textTransform: 'none',
                fontFamily: 'Roboto, arial, sans-serif',
                fontSize: '14px',
                height: '40px',
                padding: '0 12px',
                maxWidth: 'max-content',
                minWidth: 'min-content',
                '&:hover': {
                    boxShadow: '0 1px 2px 0 rgba(60, 64, 67, .30), 0 1px 3px 1px rgba(60, 64, 67, .15)',
                    backgroundColor: '#131314',
                },
                '&:disabled': {
                    backgroundColor: '#13131461',
                    borderColor: '#8e918f1f',
                    color: '#e3e3e3',
                    opacity: 0.38,
                },
            }}
            disabled={disabled}
        >
            Sign in with Google
        </Button>
    );
};

export default GoogleSignInButton;
