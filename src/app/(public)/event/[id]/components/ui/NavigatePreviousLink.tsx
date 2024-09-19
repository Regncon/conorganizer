import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Button, IconButton, Typography, useMediaQuery, type SxProps } from '@mui/material';
import Link from 'next/link';
type Props = {
    previousNavigationId: string | undefined;
};

const NavigatePreviousLink = ({ previousNavigationId }: Props) => {
    const sx: SxProps = { fontSize: 'var(--arrow-size)', width: 'fit-content', placeSelf: 'start' };
    const href = `/event/${previousNavigationId}`;

    return previousNavigationId ?
            <Button variant="outlined" color="secondary" sx={sx} component={Link} href={href}>
                <ArrowBackIcon />
                <Typography sx={{ marginLeft: '0.5rem' }}>forrige</Typography>
            </Button>
        :   null;
};

export default NavigatePreviousLink;
