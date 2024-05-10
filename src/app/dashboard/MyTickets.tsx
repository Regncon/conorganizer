import CardBase from './CardBase';

type Props = {};

const MyTickets = ({}: Props) => {
	return (
		<CardBase
			description="Trykk for og gå til mine billetter"
			img="/my-tickets.jpg"
			imgAlt="Mine billeter"
			title="Mine billetter"
		/>
	);
};

export default MyTickets;
