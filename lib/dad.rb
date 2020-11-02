# frozen_string_literal: true

# The new dad
class Dad
  @@inst = nil
  def self.get
    if @@inst.nil?
      @@inst = Dad.new
    end
    return @@inst
  end
  def request(str)
    generic = [
      'Want to go fishing?',
      'Im gonna go grill!',
      'Gotta go get some milk...',
      'Wanna throw the pigskin around?',
      'Good stuff sport!',
      'Good stuff champ!'
    ]
    if str.downcase.include?('i love you')
      return 'I dont love you'
    elsif str.downcase.include?('exit')
      return 'You will never leave'
    else
      ran = rand(generic.length)
      return generic[ran]
    end
  end
end
@dad = Dad.new

