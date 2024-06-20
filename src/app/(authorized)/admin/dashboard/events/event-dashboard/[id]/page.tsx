import EventCardBig from '$app/(public)/EventCardBig';
import { Paper } from '@mui/material';
import Link from 'next/link';
import EventDashboard from './EventDashboard';

const Page = async () => {
    return <Paper>
        <EventDashboard />
    </Paper>;
};

export default Page;
