import { Paper } from '@mui/material';
import EventDashboardTabs from '../components/EventDashboardTabs';
import Room from './components/Room';

type Props = {
    params: {
        id: string;
    };
};
const page = ({ params: { id } }: Props) => {
    return (
        <Paper elevation={0}>
            <EventDashboardTabs id={id} value={3} />
            <Room id={id} />
        </Paper>
    );
};

export default page;
