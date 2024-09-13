import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { Button, IconButton, Typography, useMediaQuery, type SxProps } from '@mui/material';
import Link from 'next/link';

type Props = {
    nextNavigationId: string | undefined;
};

const NavigateNextLink = ({ nextNavigationId }: Props) => {
    const sx: SxProps = { fontSize: 'var(--arrow-size)', width: 'fit-content', placeSelf: 'end' };
    const href = `/event/${nextNavigationId}`;

    return nextNavigationId ?
            <Button variant="outlined" color="secondary" sx={sx} component={Link} href={href} replace>
                <Typography sx={{ marginLeft: '0.5rem' }}>neste</Typography>
                <ArrowForwardIcon />
            </Button>
        :   null;
};

export default NavigateNextLink;
