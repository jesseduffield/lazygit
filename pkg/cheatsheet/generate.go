// Fix newlines breaking tables
	tooltip := strings.ReplaceAll(binding.Tooltip, "\r\n", " ")
	tooltip = strings.ReplaceAll(tooltip, "\n", " ")
	tooltip = strings.ReplaceAll(tooltip, "\t", " ")
	tooltip = strings.TrimSpace(tooltip)