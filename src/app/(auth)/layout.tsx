import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Box from '@mui/material/Box';
import { redirect } from 'next/navigation';

type Props = {
    children: React.ReactNode;
};
const layout = async ({ children }: Props) => {
    const { auth } = await getAuthorizedAuth();
    if (auth === null) {
        return <Box sx={{ marginBlockStart: '5rem' }}>{children}</Box>;
    }
    redirect('/dashboard');
};

export default layout;
