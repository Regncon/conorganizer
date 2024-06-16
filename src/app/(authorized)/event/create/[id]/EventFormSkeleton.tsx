import Skeleton from '@mui/material/Skeleton/Skeleton';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';

const EventFromSkeleton = () => {
    return (
        <Grid2 container spacing="2rem">
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={120} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={120} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={84} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={84} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={84} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={84} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={129} />
            </Grid2>
            <Grid2 xs={12} sm={4}>
                <Skeleton variant="rounded" height={220} />
            </Grid2>
            <Grid2 xs={12} sm={4}>
                <Skeleton variant="rounded" sx={{ height: { xs: 56, sm: 220 } }} />
            </Grid2>
            <Grid2 xs={12} sm={4}>
                <Skeleton variant="rounded" height={220} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={309} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={90} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={80} />
            </Grid2>
        </Grid2>
    );
};

export default EventFromSkeleton;
