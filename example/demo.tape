# Output demo.gif
Output demo.mp4
# Output demo.ascii

Set Width 1750
Set Height 1250
Set LineHeight 1.3
Set FontSize 35
Set WindowBar Colorful
Set CursorBlink false
Set PlaybackSpeed 2

Require tmux
Require telnet

Hide
  Type "tmux -f /dev/null -L test new-session -- bash" Enter
  Type "tmux split-window -h -p 51 -- bash && \" Enter
  Type "tmux split-window -h -p 99 -- bash && \" Enter
  Type "tmux send-keys -t 1 'while true; do clear; done' && \" Enter
  Type "tmux send-keys -t 1 Enter && \" Enter
  Type "tmux set status && \" Enter
  Type 'tmux setw pane-border-style "fg=0" && \' Enter
  Type 'tmux setw pane-active-border-style "fg=0" && \' Enter
  Type "clear" Enter
  Sleep 2
Show

Hide
  Type 'tmux select-pane -t 0 && clear' Enter
  Sleep 1
Show

Type "./deris" Enter
Sleep 2

# tmux bind to switch to prev pane
Ctrl+B
Type ";"

Sleep 2
Type "telnet localhost 6969" Enter
Sleep 1
Type@250ms "SET cat meow" Enter
Sleep 0.8
Type@250ms "SET cow moo" Enter
Sleep 0.8
Type@250ms "GET cow" Enter
Sleep 0.8
Type@250ms "GET dog" Enter
Sleep 0.8
Type@250ms "DEL cat" Enter
Sleep 0.8
Type@250ms "GET cat" Enter
Sleep 0.8
Type@250ms "exit" Enter
Sleep 1
Type@250ms "exit" Enter
Sleep 0.5

Hide
  Ctrl+C
  Ctrl+C
  Ctrl+C
  Type "exit" Enter
  Sleep 0.5
Show


Hide
  Ctrl+C
  Sleep 3
Show

Type "clear" Enter
Type "cat snapshot.log" Enter
Sleep 2
Type "./deris" Enter
Sleep 2
