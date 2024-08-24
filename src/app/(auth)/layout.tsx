import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import { redirect } from 'next/navigation';

type Props = {
    children: React.ReactNode;
};
const layout = async ({ children }: Props) => {
    const { auth } = await getAuthorizedAuth();
    if (auth === null) {
        return (
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateRows: '100dvh',
                    gridTemplateColumns: 'minmax(auto, 40rem)',
                    justifyContent: 'center',
                    alignItems: 'center',
                }}
            >
                <Paper sx={{ display: 'grid', placeContent: 'center', placeItems: 'center', height: '65%' }}>
                    {children}
                </Paper>
            </Box>
        );
    }
    return null;
    //redirect('/dashboard');
};

export default layout;
