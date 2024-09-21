import { Button } from '@mui/material';
import LaunchIcon from '@mui/icons-material/Launch';
type Props = {};

const BuyTicketButton = ({}: Props) => {
    return (
        <Button
            fullWidth
            variant="contained"
            href="https://event.checkin.no/73685/regncon-xxxii-2024"
            color="secondary"
        >
            Kjøp billett <LaunchIcon sx={{ marginInlineStart: '1rem' }} />
        </Button>
    );
};

export default BuyTicketButton;
