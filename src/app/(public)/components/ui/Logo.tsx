import RegnconLogo2024 from '$ui/RegnconLogo2024';
import { Box } from '@mui/material';
type Props = {};

const Logo = ({ }: Props) => {
    return (
        <Box sx={{ display: 'flex', justifyContent: 'center' }}>
            <RegnconLogo2024 spin={true} size="large" />
        </Box>
    );
};

export default Logo;
