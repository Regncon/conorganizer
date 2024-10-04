import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import { Box, Typography } from '@mui/material';

type Props = {
    over18: boolean;
};

const Over18 = ({ over18 }: Props) => {
    return (
        <Box sx={{ display: 'flex', justifyContent: 'start', alignItems: 'center' }}>
            {over18 ?
                <>
                    <AdultsOnlyIcon chipMargin={false} />
                    <Typography sx={{ paddingLeft: '0.5rem', fontWeight: 'bold' }}>Over 18</Typography>
                </>
                : <>
                    <ChildFriendlyIcon chipMargin={false} />
                    <Typography sx={{ paddingLeft: '0.5rem', fontWeight: 'bold' }}>Under 18</Typography>
                </>
            }
        </Box>
    );
};
export default Over18;
