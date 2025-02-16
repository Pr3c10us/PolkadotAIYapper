package command

import (
	"fmt"
	"math/rand"
)

const (
	USERCENTRICContent                   = "Generate 10 Polkadot ecosystem topics that directly address end-user needs and experiences. Focus on practical applications, user benefits, and real-world use cases that would matter to everyday blockchain users. Topics should be accessible to non-technical users while remaining substantive. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	SIMPLIFYWEB3JARGON                   = "Generate 10 Polkadot ecosystem topics that explain complex technical concepts in simple terms. Focus on breaking down technical jargon into accessible language while maintaining accuracy. Each topic should help bridge the gap between technical and non-technical understanding. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	BLOCKCHAININTEROPERABILITY           = "Generate 10 topics focusing on Polkadot's cross-chain communication capabilities, parachain interactions, and interoperability solutions. Include topics about cross-chain bridges, XCMP, and how Polkadot enables different blockchains to work together. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	EXPLOREWEB3GOVERNANCEDAOS            = "Generate 10 topics exploring Polkadot's governance mechanisms, including OpenGov, referenda, council operations, and treasury management. Focus on how DAOs operate within the Polkadot ecosystem and democratic decision-making processes. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	EXPLAINCOMPLEXCONCEPTS               = "Generate 10 in-depth technical topics about Polkadot's architecture, including consensus mechanisms, nominated proof-of-stake, parachain auctions, and cross-chain messaging. Focus on sophisticated concepts that would interest developers and technical users. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	NARRATIVESANDCASESTUDIES             = "Generate 10 topics centered around real-world applications, success stories, and case studies within the Polkadot ecosystem. Include specific examples of projects, partnerships, and implementations that demonstrate Polkadot's impact. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	DEVELOPERFOCUSEDCONTENT              = "Generate 10 technical topics specifically for developers building on Polkadot, including Substrate framework, smart contract development, parachain deployment, and tooling. Focus on practical development challenges and solutions. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	ADDRESSINGWEB3MYTHSANDMISCONCEPTIONS = "Generate 10 topics that address common misconceptions and myths about Polkadot and its technology. Focus on clarifying misunderstandings about scalability, security, decentralization, and other aspects of the ecosystem. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	COMMONQUESTIONS                      = "Generate 10 topics based on frequently asked questions about Polkadot, including staking, governance participation, parachain investments, and network functionality. Focus on questions that consistently arise in community discussions. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
	EXPLOREEMERGINGTRENDS                = "Generate 10 forward-looking topics about emerging trends and future developments in the Polkadot ecosystem, including upcoming protocol upgrades, new parachain launches, and potential industry impacts. Focus on innovations and future possibilities. Return the result as an array of strings formatted as [\"topic 1\", \"topic 2\", etc.]."
)

func (service *Tweet) RandomStandardPrompt() string {
	list := []string{
		USERCENTRICContent,
		SIMPLIFYWEB3JARGON,
		BLOCKCHAININTEROPERABILITY,
		EXPLOREWEB3GOVERNANCEDAOS,
		EXPLAINCOMPLEXCONCEPTS,
		NARRATIVESANDCASESTUDIES,
		DEVELOPERFOCUSEDCONTENT,
		ADDRESSINGWEB3MYTHSANDMISCONCEPTIONS,
		COMMONQUESTIONS,
		EXPLOREEMERGINGTRENDS,
	}

	randIndex := rand.Intn(len(list) - 0)
	return list[randIndex]
}

func (service *Tweet) ProductListPrompt() string {
	return "As a blockchain technology expert, provide a JSON array of exactly 10 Polkadot ecosystem projects, with a mix of:\n\nA. Established projects (4 slots) that meet these criteria:\n- Live on mainnet\n- Successful parachain slot auction history\n- Minimum $1M TVL\n- Valid security audits\n\nB. Emerging projects (3 slots) that meet these criteria:\n- Currently in testnet/beta\n- Active development (weekly commits)\n- Public roadmap\n- Secured funding/grants\n\nC. Early-stage projects (3 slots) that meet these criteria:\n- Announced within last 6 months\n- Novel use case or technology\n- Clear development timeline\n- Backing from recognized teams/VCs\n\nFormat requirements:\n- Strict JSON array format: [\"name1\", \"name2\", ...]\n- Exactly 10 elements total\n- Project names must match official branding\n- Double quotes required\n- No trailing comma\n- No whitespace between elements\n\nExample of correct formatting:\n[\"Acala\",\"Moonbeam\",\"NewProject\",\"UpcomingDapp\",\"ProjectName5\",\"ProjectName6\",\"ProjectName7\",\"ProjectName8\",\"ProjectName9\",\"ProjectName10\"]\n\nReturn only the JSON array, with no additional text, formatting, or explanation."
}

func (service *Tweet) ProductTopicPrompt(product string) string {
	return fmt.Sprintf("You are a blockchain expert. Generate exactly 10 fascinating single-sentence topics about %s within the Polkadot ecosystem.\n\n  \n\nRequirements:\n\nEach topic must explicitly mention %s and its connection to Polkadot.\n\nOnly one sentence per topic.\n\nTopics must be unique, specific, and truly engaging.\n\nExplore notable, groundbreaking, or unusual aspects of %s within Polkadot.\n\nBase all topics on real features, achievements, or verified facts.\n\nAvoid generic blockchain statements—focus on what makes %s stand out in Polkadot.\n\nFormat:\n\nReturn the output as a list of exactly 10 items in this format:\n\n[\"topic 1\", \"topic 2\", \"topic 3\", ..., \"topic 10\"]\n\nUse clear, engaging language.\n\nInclude specific details, metrics, or unique terminology when relevant.\n\nExample Output:\n\n[\"Astar Network pioneered the first 'Build2Earn' program in the Polkadot ecosystem, rewarding developers with native tokens for deploying smart contracts.\", \"Moonbeam seamlessly integrates Ethereum dApps into the Polkadot ecosystem, enabling cross-chain interoperability with Substrate-based parachains.\"]\n\nReturn only the list—no extra text.", product, product, product, product)
}
