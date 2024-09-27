import RegnconLogo2024 from '$ui/RegnconLogo2024';
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
            <RegnconLogo2024 spin={true} size="large" />
        </Box>
    );
};

export default Logo;
