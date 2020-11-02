# frozen_string_literal: true

require 'discordrb'
require_relative 'dad'
require_relative 'events'

# bot class
class Bot
  @@inst = nil
  def self.get
    if @@inst.nil?
      @@inst = Bot.new
    end
    return @@inst
  end
  @bot = nil
  attr_accessor :bot
  attr_accessor :t_file
  def initialize
    @t_file = File.new('./bot.token', 'a+')
    arg = IO.readlines('./bot.token')
    puts 'Created file'
    if !arg[0].nil?
      @bot = Discordrb::Bot.new token: arg[0]
    else
      puts 'No token!'
      sleep(3)
      exit(1)
    end
  end
end
gbot = Bot.new
gbot.bot.include! EvContainer
gbot.bot.run


