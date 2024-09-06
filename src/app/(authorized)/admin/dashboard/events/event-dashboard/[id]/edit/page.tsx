import { Paper } from '@mui/material';
import EventDashboardTabs from '../EventDashboardTabs';
import MainEvent from '$app/(public)/event/[id]/MainEvent';
import Edit from './Edit';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    return (
        <Paper elevation={0}>
            <EventDashboardTabs id={id} value={4} />
            <Edit id={id} />
        </Paper>
    );
};

export default Page;
