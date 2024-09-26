import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { Card, CardActionArea, CardContent, CardHeader, type SxProps, type Theme } from '@mui/material';
import Image from 'next/image';
import rook from '$lib/image/rook.svg';
import type { EventCardProps } from '../../../../lib/types';
import TrashButton from './TrashButton';
import { faCircleCheck } from '@fortawesome/free-solid-svg-icons/faCircleCheck';
import { faPencil } from '@fortawesome/free-solid-svg-icons/faPencil';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckDouble } from '@fortawesome/free-solid-svg-icons';
import { getParticipantByUser } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';

export default function EventCardBig({
    title,
    gameMaster,
    shortDescription,
    system,
    backgroundImage = '/dice-small.webp',
    myEventBar = false,
    myEventBarSubmitted = false,
    myEventDocId,
    isAccepted = false,
}: EventCardProps) {
    const circleCheckOrPencilIcon =
        isAccepted ? faCheckDouble
            : myEventBarSubmitted ? faCircleCheck
                : faPencil;
    const SuccessOrWarningColor =
        isAccepted ? 'success.dark'
            : myEventBarSubmitted ? 'success.light'
                : 'warning.main';

    return (
        <Card
            sx={{
                backgroundImage: `url(${backgroundImage ? backgroundImage : '/dice-small.webp'})`,
                minHeight: `${myEventBar ? 'calc(267px + 65px)' : '267px'}`,
                maxHeight: `${myEventBar ? 'calc(267px + 65px)' : '267px'}`,
                backgroundSize: 'cover',
                borderRadius: '1.75rem',
                display: 'grid',
                minWidth: '306px',
                opacity: isAccepted ? '0.5' : '1',
            }}
        >
            <CardActionArea
                disabled={isAccepted}
                sx={{
                    ['&:has(.disable-ripple) > .MuiTouchRipple-root']: {
                        display: 'none',
                    },
                    display: 'grid',
                    gridTemplateRows: `${myEventBar ? 'auto' : ''} 1fr 0.895fr`,
                    placeItems: 'end',
                    gridTemplateColumns: 'subgrid',
                }}
            >
                {myEventBar ?
                    <>
                        <Box
                            sx={{
                                display: 'grid',
                                gridAutoFlow: 'column',
                                placeContent: 'space-between',
                                width: '100%',
                                padding: '1rem',
                            }}
                        >
                            <Box sx={{ display: 'flex', gap: '0.5rem', color: 'success.main', placeItems: 'center' }}>
                                <Typography component="span" sx={{ color: SuccessOrWarningColor }}>
                                    <FontAwesomeIcon icon={circleCheckOrPencilIcon} size="2x" />
                                </Typography>
                                <Typography sx={{ color: SuccessOrWarningColor }}>
                                    {isAccepted ?
                                        'Godtatt'
                                        : myEventBarSubmitted ?
                                            'Sendt inn'
                                            : 'Kladd'}
                                </Typography>
                            </Box>
                            <TrashButton docId={myEventDocId} />
                        </Box>
                    </>
                    : null}

                <CardHeader
                    title={title}
                    titleTypographyProps={{ fontSize: '1.8rem' }}
                    sx={{
                        alignItems: 'flex-end',
                        marginBlockEnd: '1rem',
                        wordBreak: 'break-all',
                        overflow: 'clip',
                        maxHeight: '7rem',
                        padding: '0',
                        paddingInlineStart: '1rem',
                        placeSelf: 'end start',
                        '.MuiCardHeader-content': {
                            maxHeight: '7rem',
                        },
                    }}
                />
                <CardContent
                    sx={{
                        backgroundColor: 'rgba(0,0,0,0.5)',
                        backdropFilter: 'blur(4px)',
                        padding: '1rem',
                        height: '100%',
                        width: '100%',
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
                        <Typography fontSize="1rem"> {system} </Typography>
                        <Box sx={{ display: 'flex', gap: '1rem' }}>
                            {/*   <Box component={Image} priority src={rook} alt="rook icon" />
                            <Box component={Image} priority src={rook} alt="rook icon" />
                            <Box component={Image} priority src={rook} alt="rook icon" />
                            <Box component={Image} priority src={rook} alt="rook icon" />
                            */}
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
                            fontSize: '1rem',
                            width: '90%',
                        }}
                    >
                        {shortDescription}
                    </Typography>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}
