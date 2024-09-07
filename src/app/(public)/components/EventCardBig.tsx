import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { Card, CardActionArea, CardContent, CardHeader, type SxProps, type Theme } from '@mui/material';
import Image from 'next/image';
import rook from '$lib/image/rook.svg';
import type { EventCardProps } from '../../../lib/types';
import TrashButton from './TrashButton';
import { faCircleCheck } from '@fortawesome/free-solid-svg-icons/faCircleCheck';
import { faPencil } from '@fortawesome/free-solid-svg-icons/faPencil';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export default function EventCardBig({
    title,
    gameMaster,
    shortDescription,
    system,
    backgroundImage = 'blekksprut2.jpg',
    myEventBar = false,
    myEventBarSubmitted = false,
    myEventDocId,
}: EventCardProps) {
    const circleCheckOrPencilIcon = myEventBarSubmitted ? faCircleCheck : faPencil;
    const SuccessOrWarningColor = myEventBarSubmitted ? 'success.main' : 'warning.main';
    const width: SxProps<Theme> = {
        width: '100vw',
        maxWidth: '430px',
    };
    return (
        <Card
            sx={{
                backgroundImage: `url(/${backgroundImage})`,
                ...width,
                height: '269px',
                backgroundSize: 'cover',
                borderRadius: '1.75rem',
            }}
        >
            <CardActionArea
                sx={{
                    ['&:has(.disable-ripple) > .MuiTouchRipple-root']: {
                        display: 'none',
                    },
                    display: 'grid',
                    gridTemplateRows: '1fr 1fr',
                    placeItems: 'end',
                }}
            >
                <Box
                    sx={{
                        display: 'flex',
                        placeSelf: 'start',
                        placeContent: 'space-between',
                        padding: '1rem',
                        width: '100%',
                    }}
                >
                    {myEventBar ?
                        <>
                            <Box sx={{ display: 'flex', gap: '0.5rem', color: 'success.main', placeItems: 'center' }}>
                                <Typography component="span" sx={{ color: SuccessOrWarningColor }}>
                                    <FontAwesomeIcon icon={circleCheckOrPencilIcon} size="2x" />
                                </Typography>
                                <Typography sx={{ color: SuccessOrWarningColor }}>
                                    {myEventBarSubmitted ? 'Sendt inn' : 'Kladd'}
                                </Typography>
                            </Box>
                            <TrashButton docId={myEventDocId} />
                        </>
                    :   null}
                </Box>
                <CardHeader
                    title={title}
                    titleTypographyProps={{ fontSize: '1.8rem' }}
                    sx={{
                        ...width,
                        height: '72px',
                        alignItems: 'flex-end',
                        padding: '1rem',
                        wordBreak: 'break-all',
                        display: '-webkit-box',
                        overflow: 'clip',
                        WebkitLineClamp: '2',
                        WebkitBoxOrient: 'vertical',
                        placeSelf: 'end start',
                    }}
                />
                <CardContent
                    sx={{
                        height: '126px',
                        backgroundColor: 'rgba(0,0,0,0.5)',
                        backdropFilter: 'blur(4px)',
                        padding: '1rem',
                        ...width,
                    }}
                >
                    <Typography sx={{ fontWeight: 'bold', fontSize: '1.1rem' }}> {gameMaster} </Typography>
                    <Box
                        sx={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            color: 'primary.main',
                            paddingBottom: '0.5rem',
                        }}
                    >
                        <Typography> {system} </Typography>
                        <Box sx={{ display: 'flex', gap: '1rem' }}>
                            <Box component={Image} priority src={rook} alt="rook icon" />
                            <Box component={Image} priority src={rook} alt="rook icon" />
                            <Box component={Image} priority src={rook} alt="rook icon" />
                            <Box component={Image} priority src={rook} alt="rook icon" />
                        </Box>
                    </Box>
                    <Typography
                        sx={{
                            color: 'white',
                            wordBreak: 'break-all',
                            display: '-webkit-box',
                            overflow: 'clip',
                            WebkitLineClamp: '2',
                            WebkitBoxOrient: 'vertical',
                        }}
                    >
                        {shortDescription}
                    </Typography>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}