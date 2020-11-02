require 'discordrb'
module EvContainer
  extend Discordrb::EventContainer
  message do |event|
    if event.content.start_with?('<@!772951738247544882>')
      str = Dad.get.request event.content
      event.respond str
    end
  end
end
