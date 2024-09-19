import { Box } from '@mui/material';
import Image from 'next/image';
type Props = {};

const Logo = ({ }: Props) => {
    return (
        <Box
            sx={{
                maxWidth: '430px',
                maxHeight: '430px',
                margin: 'auto',
                width: '100vw',
                aspectRatio: '1/1',
                marginBlockStart: '0.5rem',
                marginBlockEnd: '1rem',
                position: 'relative',
            }}
        >
            <Image src="/BetaTestRegnconLogo.webp" fill alt="logo" priority={true} sizes="100vw" />
        </Box>
    );
};

export default Logo;
