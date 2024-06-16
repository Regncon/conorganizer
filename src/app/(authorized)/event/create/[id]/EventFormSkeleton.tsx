import Skeleton from '@mui/material/Skeleton/Skeleton';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';

const EventFromSkeleton = () => {
    const skeletonHeight = 53;
    return (
        <Grid2 container spacing="2rem">
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={skeletonHeight} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={skeletonHeight} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={skeletonHeight} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={skeletonHeight} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={skeletonHeight} />
            </Grid2>
            <Grid2 xs={12} sm={6} md={3}>
                <Skeleton variant="rounded" height={skeletonHeight} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={129} />
            </Grid2>
            <Grid2 xs={12} sm={4}>
                <Skeleton variant="rounded" height={220} />
            </Grid2>
            <Grid2 xs={12} sm={4}>
                <Skeleton variant="rounded" sx={{ height: { xs: skeletonHeight, sm: '220px' } }} />
            </Grid2>
            <Grid2 xs={12} sm={4}>
                <Skeleton variant="rounded" height={220} />
            </Grid2>
            <Grid2 xs={12}>
                <Skeleton variant="rounded" height={380} />
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
