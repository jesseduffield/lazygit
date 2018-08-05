package termbox

import "math"
import "syscall"
import "unsafe"
import "unicode/utf16"
import "github.com/mattn/go-runewidth"

type (
	wchar     uint16
	short     int16
	dword     uint32
	word      uint16
	char_info struct {
		char wchar
		attr word
	}
	coord struct {
		x short
		y short
	}
	small_rect struct {
		left   short
		top    short
		right  short
		bottom short
	}
	console_screen_buffer_info struct {
		size                coord
		cursor_position     coord
		attributes          word
		window              small_rect
		maximum_window_size coord
	}
	console_cursor_info struct {
		size    dword
		visible int32
	}
	input_record struct {
		event_type word
		_          [2]byte
		event      [16]byte
	}
	key_event_record struct {
		key_down          int32
		repeat_count      word
		virtual_key_code  word
		virtual_scan_code word
		unicode_char      wchar
		control_key_state dword
	}
	window_buffer_size_record struct {
		size coord
	}
	mouse_event_record struct {
		mouse_pos         coord
		button_state      dword
		control_key_state dword
		event_flags       dword
	}
	console_font_info struct {
		font      uint32
		font_size coord
	}
)

const (
	mouse_lmb = 0x1
	mouse_rmb = 0x2
	mouse_mmb = 0x4 | 0x8 | 0x10
	SM_CXMIN  = 28
	SM_CYMIN  = 29
)

func (this coord) uintptr() uintptr {
	return uintptr(*(*int32)(unsafe.Pointer(&this)))
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var moduser32 = syscall.NewLazyDLL("user32.dll")
var is_cjk = runewidth.IsEastAsian()

var (
	proc_set_console_active_screen_buffer = kernel32.NewProc("SetConsoleActiveScreenBuffer")
	proc_set_console_screen_buffer_size   = kernel32.NewProc("SetConsoleScreenBufferSize")
	proc_create_console_screen_buffer     = kernel32.NewProc("CreateConsoleScreenBuffer")
	proc_get_console_screen_buffer_info   = kernel32.NewProc("GetConsoleScreenBufferInfo")
	proc_write_console_output             = kernel32.NewProc("WriteConsoleOutputW")
	proc_write_console_output_character   = kernel32.NewProc("WriteConsoleOutputCharacterW")
	proc_write_console_output_attribute   = kernel32.NewProc("WriteConsoleOutputAttribute")
	proc_set_console_cursor_info          = kernel32.NewProc("SetConsoleCursorInfo")
	proc_set_console_cursor_position      = kernel32.NewProc("SetConsoleCursorPosition")
	proc_get_console_cursor_info          = kernel32.NewProc("GetConsoleCursorInfo")
	proc_read_console_input               = kernel32.NewProc("ReadConsoleInputW")
	proc_get_console_mode                 = kernel32.NewProc("GetConsoleMode")
	proc_set_console_mode                 = kernel32.NewProc("SetConsoleMode")
	proc_fill_console_output_character    = kernel32.NewProc("FillConsoleOutputCharacterW")
	proc_fill_console_output_attribute    = kernel32.NewProc("FillConsoleOutputAttribute")
	proc_create_event                     = kernel32.NewProc("CreateEventW")
	proc_wait_for_multiple_objects        = kernel32.NewProc("WaitForMultipleObjects")
	proc_set_event                        = kernel32.NewProc("SetEvent")
	proc_get_current_console_font         = kernel32.NewProc("GetCurrentConsoleFont")
	get_system_metrics                    = moduser32.NewProc("GetSystemMetrics")
)

func set_console_active_screen_buffer(h syscall.Handle) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_active_screen_buffer.Addr(),
		1, uintptr(h), 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_screen_buffer_size(h syscall.Handle, size coord) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_screen_buffer_size.Addr(),
		2, uintptr(h), size.uintptr(), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func create_console_screen_buffer() (h syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall6(proc_create_console_screen_buffer.Addr(),
		5, uintptr(generic_read|generic_write), 0, 0, console_textmode_buffer, 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return syscall.Handle(r0), err
}

func get_console_screen_buffer_info(h syscall.Handle, info *console_screen_buffer_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_screen_buffer_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func write_console_output(h syscall.Handle, chars []char_info, dst small_rect) (err error) {
	tmp_coord = coord{dst.right - dst.left + 1, dst.bottom - dst.top + 1}
	tmp_rect = dst
	r0, _, e1 := syscall.Syscall6(proc_write_console_output.Addr(),
		5, uintptr(h), uintptr(unsafe.Pointer(&chars[0])), tmp_coord.uintptr(),
		tmp_coord0.uintptr(), uintptr(unsafe.Pointer(&tmp_rect)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func write_console_output_character(h syscall.Handle, chars []wchar, pos coord) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_write_console_output_character.Addr(),
		5, uintptr(h), uintptr(unsafe.Pointer(&chars[0])), uintptr(len(chars)),
		pos.uintptr(), uintptr(unsafe.Pointer(&tmp_arg)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func write_console_output_attribute(h syscall.Handle, attrs []word, pos coord) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_write_console_output_attribute.Addr(),
		5, uintptr(h), uintptr(unsafe.Pointer(&attrs[0])), uintptr(len(attrs)),
		pos.uintptr(), uintptr(unsafe.Pointer(&tmp_arg)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_cursor_info(h syscall.Handle, info *console_cursor_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_cursor_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func get_console_cursor_info(h syscall.Handle, info *console_cursor_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_cursor_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_cursor_position(h syscall.Handle, pos coord) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_cursor_position.Addr(),
		2, uintptr(h), pos.uintptr(), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func read_console_input(h syscall.Handle, record *input_record) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_read_console_input.Addr(),
		4, uintptr(h), uintptr(unsafe.Pointer(record)), 1, uintptr(unsafe.Pointer(&tmp_arg)), 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func get_console_mode(h syscall.Handle, mode *dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_mode.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(mode)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_mode(h syscall.Handle, mode dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_mode.Addr(),
		2, uintptr(h), uintptr(mode), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func fill_console_output_character(h syscall.Handle, char wchar, n int) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_fill_console_output_character.Addr(),
		5, uintptr(h), uintptr(char), uintptr(n), tmp_coord.uintptr(),
		uintptr(unsafe.Pointer(&tmp_arg)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func fill_console_output_attribute(h syscall.Handle, attr word, n int) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_fill_console_output_attribute.Addr(),
		5, uintptr(h), uintptr(attr), uintptr(n), tmp_coord.uintptr(),
		uintptr(unsafe.Pointer(&tmp_arg)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func create_event() (out syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall6(proc_create_event.Addr(),
		4, 0, 0, 0, 0, 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return syscall.Handle(r0), err
}

func wait_for_multiple_objects(objects []syscall.Handle) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_wait_for_multiple_objects.Addr(),
		4, uintptr(len(objects)), uintptr(unsafe.Pointer(&objects[0])),
		0, 0xFFFFFFFF, 0, 0)
	if uint32(r0) == 0xFFFFFFFF {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_event(ev syscall.Handle) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_event.Addr(),
		1, uintptr(ev), 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func get_current_console_font(h syscall.Handle, info *console_font_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_current_console_font.Addr(),
		3, uintptr(h), 0, uintptr(unsafe.Pointer(info)))
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

type diff_msg struct {
	pos   short
	lines short
	chars []char_info
}

type input_event struct {
	event Event
	err   error
}

var (
	orig_cursor_info console_cursor_info
	orig_size        coord
	orig_mode        dword
	orig_screen      syscall.Handle
	back_buffer      cellbuf
	front_buffer     cellbuf
	term_size        coord
	input_mode       = InputEsc
	cursor_x         = cursor_hidden
	cursor_y         = cursor_hidden
	foreground       = ColorDefault
	background       = ColorDefault
	in               syscall.Handle
	out              syscall.Handle
	interrupt        syscall.Handle
	charbuf          []char_info
	diffbuf          []diff_msg
	beg_x            = -1
	beg_y            = -1
	beg_i            = -1
	input_comm       = make(chan Event)
	interrupt_comm   = make(chan struct{})
	cancel_comm      = make(chan bool, 1)
	cancel_done_comm = make(chan bool)
	alt_mode_esc     = false

	// these ones just to prevent heap allocs at all costs
	tmp_info   console_screen_buffer_info
	tmp_arg    dword
	tmp_coord0 = coord{0, 0}
	tmp_coord  = coord{0, 0}
	tmp_rect   = small_rect{0, 0, 0, 0}
	tmp_finfo  console_font_info
)

func get_cursor_position(out syscall.Handle) coord {
	err := get_console_screen_buffer_info(out, &tmp_info)
	if err != nil {
		panic(err)
	}
	return tmp_info.cursor_position
}

func get_term_size(out syscall.Handle) coord {
	err := get_console_screen_buffer_info(out, &tmp_info)
	if err != nil {
		panic(err)
	}
	return tmp_info.size
}

func get_win_min_size(out syscall.Handle) coord {
	x, _, err := get_system_metrics.Call(SM_CXMIN)
	y, _, err := get_system_metrics.Call(SM_CYMIN)

	if x == 0 || y == 0 {
		if err != nil {
			panic(err)
		}
	}

	err1 := get_current_console_font(out, &tmp_finfo)
	if err1 != nil {
		panic(err1)
	}

	return coord{
		x: short(math.Ceil(float64(x) / float64(tmp_finfo.font_size.x))),
		y: short(math.Ceil(float64(y) / float64(tmp_finfo.font_size.y))),
	}
}

func get_win_size(out syscall.Handle) coord {
	err := get_console_screen_buffer_info(out, &tmp_info)
	if err != nil {
		panic(err)
	}

	min_size := get_win_min_size(out)

	size := coord{
		x: tmp_info.window.right - tmp_info.window.left + 1,
		y: tmp_info.window.bottom - tmp_info.window.top + 1,
	}

	if size.x < min_size.x {
		size.x = min_size.x
	}

	if size.y < min_size.y {
		size.y = min_size.y
	}

	return size
}

func update_size_maybe() {
	size := get_win_size(out)
	if size.x != term_size.x || size.y != term_size.y {
		set_console_screen_buffer_size(out, size)
		term_size = size
		back_buffer.resize(int(size.x), int(size.y))
		front_buffer.resize(int(size.x), int(size.y))
		front_buffer.clear()
		clear()

		area := int(size.x) * int(size.y)
		if cap(charbuf) < area {
			charbuf = make([]char_info, 0, area)
		}
	}
}

var color_table_bg = []word{
	0, // default (black)
	0, // black
	background_red,
	background_green,
	background_red | background_green, // yellow
	background_blue,
	background_red | background_blue,                    // magenta
	background_green | background_blue,                  // cyan
	background_red | background_blue | background_green, // white
}

var color_table_fg = []word{
	foreground_red | foreground_blue | foreground_green, // default (white)
	0,
	foreground_red,
	foreground_green,
	foreground_red | foreground_green, // yellow
	foreground_blue,
	foreground_red | foreground_blue,                    // magenta
	foreground_green | foreground_blue,                  // cyan
	foreground_red | foreground_blue | foreground_green, // white
}

const (
	replacement_char = '\uFFFD'
	max_rune         = '\U0010FFFF'
	surr1            = 0xd800
	surr2            = 0xdc00
	surr3            = 0xe000
	surr_self        = 0x10000
)

func append_diff_line(y int) int {
	n := 0
	for x := 0; x < front_buffer.width; {
		cell_offset := y*front_buffer.width + x
		back := &back_buffer.cells[cell_offset]
		front := &front_buffer.cells[cell_offset]
		attr, char := cell_to_char_info(*back)
		charbuf = append(charbuf, char_info{attr: attr, char: char[0]})
		*front = *back
		n++
		w := runewidth.RuneWidth(back.Ch)
		if w == 0 || w == 2 && runewidth.IsAmbiguousWidth(back.Ch) {
			w = 1
		}
		x += w
		// If not CJK, fill trailing space with whitespace
		if !is_cjk && w == 2 {
			charbuf = append(charbuf, char_info{attr: attr, char: ' '})
		}
	}
	return n
}

// compares 'back_buffer' with 'front_buffer' and prepares all changes in the form of
// 'diff_msg's in the 'diff_buf'
func prepare_diff_messages() {
	// clear buffers
	diffbuf = diffbuf[:0]
	charbuf = charbuf[:0]

	var diff diff_msg
	gbeg := 0
	for y := 0; y < front_buffer.height; y++ {
		same := true
		line_offset := y * front_buffer.width
		for x := 0; x < front_buffer.width; x++ {
			cell_offset := line_offset + x
			back := &back_buffer.cells[cell_offset]
			front := &front_buffer.cells[cell_offset]
			if *back != *front {
				same = false
				break
			}
		}
		if same && diff.lines > 0 {
			diffbuf = append(diffbuf, diff)
			diff = diff_msg{}
		}
		if !same {
			beg := len(charbuf)
			end := beg + append_diff_line(y)
			if diff.lines == 0 {
				diff.pos = short(y)
				gbeg = beg
			}
			diff.lines++
			diff.chars = charbuf[gbeg:end]
		}
	}
	if diff.lines > 0 {
		diffbuf = append(diffbuf, diff)
		diff = diff_msg{}
	}
}

func get_ct(table []word, idx int) word {
	idx = idx & 0x0F
	if idx >= len(table) {
		idx = len(table) - 1
	}
	return table[idx]
}

func cell_to_char_info(c Cell) (attr word, wc [2]wchar) {
	attr = get_ct(color_table_fg, int(c.Fg)) | get_ct(color_table_bg, int(c.Bg))
	if c.Fg&AttrReverse|c.Bg&AttrReverse != 0 {
		attr = (attr&0xF0)>>4 | (attr&0x0F)<<4
	}
	if c.Fg&AttrBold != 0 {
		attr |= foreground_intensity
	}
	if c.Bg&AttrBold != 0 {
		attr |= background_intensity
	}

	r0, r1 := utf16.EncodeRune(c.Ch)
	if r0 == 0xFFFD {
		wc[0] = wchar(c.Ch)
		wc[1] = ' '
	} else {
		wc[0] = wchar(r0)
		wc[1] = wchar(r1)
	}
	return
}

func move_cursor(x, y int) {
	err := set_console_cursor_position(out, coord{short(x), short(y)})
	if err != nil {
		panic(err)
	}
}

func show_cursor(visible bool) {
	var v int32
	if visible {
		v = 1
	}

	var info console_cursor_info
	info.size = 100
	info.visible = v
	err := set_console_cursor_info(out, &info)
	if err != nil {
		panic(err)
	}
}

func clear() {
	var err error
	attr, char := cell_to_char_info(Cell{
		' ',
		foreground,
		background,
	})

	area := int(term_size.x) * int(term_size.y)
	err = fill_console_output_attribute(out, attr, area)
	if err != nil {
		panic(err)
	}
	err = fill_console_output_character(out, char[0], area)
	if err != nil {
		panic(err)
	}
	if !is_cursor_hidden(cursor_x, cursor_y) {
		move_cursor(cursor_x, cursor_y)
	}
}

func key_event_record_to_event(r *key_event_record) (Event, bool) {
	if r.key_down == 0 {
		return Event{}, false
	}

	e := Event{Type: EventKey}
	if input_mode&InputAlt != 0 {
		if alt_mode_esc {
			e.Mod = ModAlt
			alt_mode_esc = false
		}
		if r.control_key_state&(left_alt_pressed|right_alt_pressed) != 0 {
			e.Mod = ModAlt
		}
	}

	ctrlpressed := r.control_key_state&(left_ctrl_pressed|right_ctrl_pressed) != 0

	if r.virtual_key_code >= vk_f1 && r.virtual_key_code <= vk_f12 {
		switch r.virtual_key_code {
		case vk_f1:
			e.Key = KeyF1
		case vk_f2:
			e.Key = KeyF2
		case vk_f3:
			e.Key = KeyF3
		case vk_f4:
			e.Key = KeyF4
		case vk_f5:
			e.Key = KeyF5
		case vk_f6:
			e.Key = KeyF6
		case vk_f7:
			e.Key = KeyF7
		case vk_f8:
			e.Key = KeyF8
		case vk_f9:
			e.Key = KeyF9
		case vk_f10:
			e.Key = KeyF10
		case vk_f11:
			e.Key = KeyF11
		case vk_f12:
			e.Key = KeyF12
		default:
			panic("unreachable")
		}

		return e, true
	}

	if r.virtual_key_code <= vk_delete {
		switch r.virtual_key_code {
		case vk_insert:
			e.Key = KeyInsert
		case vk_delete:
			e.Key = KeyDelete
		case vk_home:
			e.Key = KeyHome
		case vk_end:
			e.Key = KeyEnd
		case vk_pgup:
			e.Key = KeyPgup
		case vk_pgdn:
			e.Key = KeyPgdn
		case vk_arrow_up:
			e.Key = KeyArrowUp
		case vk_arrow_down:
			e.Key = KeyArrowDown
		case vk_arrow_left:
			e.Key = KeyArrowLeft
		case vk_arrow_right:
			e.Key = KeyArrowRight
		case vk_backspace:
			if ctrlpressed {
				e.Key = KeyBackspace2
			} else {
				e.Key = KeyBackspace
			}
		case vk_tab:
			e.Key = KeyTab
		case vk_enter:
			e.Key = KeyEnter
		case vk_esc:
			switch {
			case input_mode&InputEsc != 0:
				e.Key = KeyEsc
			case input_mode&InputAlt != 0:
				alt_mode_esc = true
				return Event{}, false
			}
		case vk_space:
			if ctrlpressed {
				// manual return here, because KeyCtrlSpace is zero
				e.Key = KeyCtrlSpace
				return e, true
			} else {
				e.Key = KeySpace
			}
		}

		if e.Key != 0 {
			return e, true
		}
	}

	if ctrlpressed {
		if Key(r.unicode_char) >= KeyCtrlA && Key(r.unicode_char) <= KeyCtrlRsqBracket {
			e.Key = Key(r.unicode_char)
			if input_mode&InputAlt != 0 && e.Key == KeyEsc {
				alt_mode_esc = true
				return Event{}, false
			}
			return e, true
		}
		switch r.virtual_key_code {
		case 192, 50:
			// manual return here, because KeyCtrl2 is zero
			e.Key = KeyCtrl2
			return e, true
		case 51:
			if input_mode&InputAlt != 0 {
				alt_mode_esc = true
				return Event{}, false
			}
			e.Key = KeyCtrl3
		case 52:
			e.Key = KeyCtrl4
		case 53:
			e.Key = KeyCtrl5
		case 54:
			e.Key = KeyCtrl6
		case 189, 191, 55:
			e.Key = KeyCtrl7
		case 8, 56:
			e.Key = KeyCtrl8
		}

		if e.Key != 0 {
			return e, true
		}
	}

	if r.unicode_char != 0 {
		e.Ch = rune(r.unicode_char)
		return e, true
	}

	return Event{}, false
}

func input_event_producer() {
	var r input_record
	var err error
	var last_button Key
	var last_button_pressed Key
	var last_state = dword(0)
	var last_x, last_y = -1, -1
	handles := []syscall.Handle{in, interrupt}
	for {
		err = wait_for_multiple_objects(handles)
		if err != nil {
			input_comm <- Event{Type: EventError, Err: err}
		}

		select {
		case <-cancel_comm:
			cancel_done_comm <- true
			return
		default:
		}

		err = read_console_input(in, &r)
		if err != nil {
			input_comm <- Event{Type: EventError, Err: err}
		}

		switch r.event_type {
		case key_event:
			kr := (*key_event_record)(unsafe.Pointer(&r.event))
			ev, ok := key_event_record_to_event(kr)
			if ok {
				for i := 0; i < int(kr.repeat_count); i++ {
					input_comm <- ev
				}
			}
		case window_buffer_size_event:
			sr := *(*window_buffer_size_record)(unsafe.Pointer(&r.event))
			input_comm <- Event{
				Type:   EventResize,
				Width:  int(sr.size.x),
				Height: int(sr.size.y),
			}
		case mouse_event:
			mr := *(*mouse_event_record)(unsafe.Pointer(&r.event))
			ev := Event{Type: EventMouse}
			switch mr.event_flags {
			case 0, 2:
				// single or double click
				cur_state := mr.button_state
				switch {
				case last_state&mouse_lmb == 0 && cur_state&mouse_lmb != 0:
					last_button = MouseLeft
					last_button_pressed = last_button
				case last_state&mouse_rmb == 0 && cur_state&mouse_rmb != 0:
					last_button = MouseRight
					last_button_pressed = last_button
				case last_state&mouse_mmb == 0 && cur_state&mouse_mmb != 0:
					last_button = MouseMiddle
					last_button_pressed = last_button
				case last_state&mouse_lmb != 0 && cur_state&mouse_lmb == 0:
					last_button = MouseRelease
				case last_state&mouse_rmb != 0 && cur_state&mouse_rmb == 0:
					last_button = MouseRelease
				case last_state&mouse_mmb != 0 && cur_state&mouse_mmb == 0:
					last_button = MouseRelease
				default:
					last_state = cur_state
					continue
				}
				last_state = cur_state
				ev.Key = last_button
				last_x, last_y = int(mr.mouse_pos.x), int(mr.mouse_pos.y)
				ev.MouseX = last_x
				ev.MouseY = last_y
			case 1:
				// mouse motion
				x, y := int(mr.mouse_pos.x), int(mr.mouse_pos.y)
				if last_state != 0 && (last_x != x || last_y != y) {
					ev.Key = last_button_pressed
					ev.Mod = ModMotion
					ev.MouseX = x
					ev.MouseY = y
					last_x, last_y = x, y
				} else {
					ev.Type = EventNone
				}
			case 4:
				// mouse wheel
				n := int16(mr.button_state >> 16)
				if n > 0 {
					ev.Key = MouseWheelUp
				} else {
					ev.Key = MouseWheelDown
				}
				last_x, last_y = int(mr.mouse_pos.x), int(mr.mouse_pos.y)
				ev.MouseX = last_x
				ev.MouseY = last_y
			default:
				ev.Type = EventNone
			}
			if ev.Type != EventNone {
				input_comm <- ev
			}
		}
	}
}
