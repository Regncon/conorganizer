import { Button } from '@mui/material';
import Link from 'next/link';

type Props = {
    parentEventId: string;
    isAdmin: boolean;
};

const GoToEventAdministrationButton = ({ parentEventId, isAdmin }: Props) => {
    return isAdmin ?
            <Button
                fullWidth
                variant="contained"
                component={Link}
                href={`/admin/dashboard/events/event-dashboard/${parentEventId}/edit`}
                prefetch
                sx={{ textDecoration: 'none', placeSelf: 'center', gridColumn: 'span 2' }}
            >
                Administrer arrangementet
            </Button>
        :   null;
};

export default GoToEventAdministrationButton;
