import { Paper } from '@mui/material';
import EventDashboardTabs from '../components/EventDashboardTabs';
import Edit from './components/Edit';

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
