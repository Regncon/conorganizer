import { Paper } from '@mui/material';
import EventDashboardTabs from '../EventDashboardTabs';
import Edit from './Edit';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    return (
        <Paper>
            <EventDashboardTabs id={id} value={4} />
            <Edit id={id} />
        </Paper>
    );
};

export default Page;
