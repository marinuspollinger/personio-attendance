complete -c personio-cli -s h -l help -d 'Show this help'
complete -c personio-cli -l status -d 'Get Current Times'
complete -c personio-cli -l break-end-time -d 'Set Break End Time (2006-01-02 15:04)' -r
complete -c personio-cli -l break-start-time -d 'Set Break Start Time (2006-01-02 15:04)' -r
complete -c personio-cli -l start-break -d 'Start the Break now!'
complete -c personio-cli -l end-break -d 'End the Break now!'
complete -c personio-cli -l send -d 'Send Current Times to Personio API'
complete -c personio-cli -l yes -d 'Immediatly write to personio, without checking times first'
