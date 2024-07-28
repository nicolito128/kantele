package gateway

// Intents is an extension of the Bit structure used when identifying with discord
type Intents int64

// Constants for the different bit offsets of Intents
const (
	IntentGuilds Intents = 1 << iota
	IntentGuildMembers
	IntentGuildModeration
	IntentGuildEmojisAndStickers
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping
	IntentMessageContent
	IntentGuildScheduledEvents
	_
	_
	_
	IntentAutoModerationConfiguration
	IntentAutoModerationExecution
	_
	_
	IntentGuildMessagePolls
	IntentDirectMessagePolls

	IntentsGuild = IntentGuilds |
		IntentGuildMembers |
		IntentGuildModeration |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildPresences |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentGuildScheduledEvents |
		IntentGuildMessagePolls

	IntentsDirectMessage = IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentDirectMessagePolls

	IntentsMessagePolls = IntentGuildMessagePolls |
		IntentDirectMessagePolls

	IntentsNonPrivileged = IntentGuilds |
		IntentGuildModeration |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution |
		IntentGuildMessagePolls |
		IntentDirectMessagePolls

	IntentsPrivileged = IntentGuildMembers |
		IntentGuildPresences | IntentMessageContent

	IntentsAll = IntentsNonPrivileged |
		IntentsPrivileged

	IntentsDefault = IntentsNone

	IntentsNone Intents = 0
)
