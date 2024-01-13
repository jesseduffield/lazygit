package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// GDScript lexer.
var GDScript = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "GDScript",
		Aliases:   []string{"gdscript", "gd"},
		Filenames: []string{"*.gd"},
		MimeTypes: []string{"text/x-gdscript", "application/x-gdscript"},
	},
	gdscriptRules,
))

func gdscriptRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`^(\s*)([rRuUbB]{,2})("""(?:.|\n)*?""")`, ByGroups(Text, LiteralStringAffix, LiteralStringDoc), nil},
			{`^(\s*)([rRuUbB]{,2})('''(?:.|\n)*?''')`, ByGroups(Text, LiteralStringAffix, LiteralStringDoc), nil},
			{`[^\S\n]+`, Text, nil},
			{`#.*$`, CommentSingle, nil},
			{`[]{}:(),;[]`, Punctuation, nil},
			{`\\\n`, Text, nil},
			{`\\`, Text, nil},
			{`(in|and|or|not)\b`, OperatorWord, nil},
			{`!=|==|<<|>>|&&|\+=|-=|\*=|/=|%=|&=|\|=|\|\||[-~+/*%=<>&^.!|$]`, Operator, nil},
			Include("keywords"),
			{`(def)((?:\s|\\\s)+)`, ByGroups(Keyword, Text), Push("funcname")},
			{`(class)((?:\s|\\\s)+)`, ByGroups(Keyword, Text), Push("classname")},
			Include("builtins"),
			{`([rR]|[uUbB][rR]|[rR][uUbB])(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("tdqs")},
			{`([rR]|[uUbB][rR]|[rR][uUbB])(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("tsqs")},
			{`([rR]|[uUbB][rR]|[rR][uUbB])(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("dqs")},
			{`([rR]|[uUbB][rR]|[rR][uUbB])(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("sqs")},
			{`([uUbB]?)(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("stringescape", "tdqs")},
			{`([uUbB]?)(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("stringescape", "tsqs")},
			{`([uUbB]?)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("stringescape", "dqs")},
			{`([uUbB]?)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("stringescape", "sqs")},
			Include("name"),
			Include("numbers"),
		},
		"keywords": {
			{Words(``, `\b`,
				`if`, `elif`, `else`, `for`, `do`,
				`while`, `switch`, `case`, `break`, `continue`,
				`pass`, `return`, `class`, `extends`, `tool`,
				`signal`, `func`, `static`, `const`, `enum`,
				`var`, `onready`, `export`, `setget`, `breakpoint`), Keyword, nil},
		},
		"builtins": {
			{Words(`(?<!\.)`, `\b`,
				`Color8`, `ColorN`, `abs`, `acos`, `asin`,
				`assert`, `atan`, `atan2`, `bytes2var`, `ceil`,
				`clamp`, `convert`, `cos`, `cosh`, `db2linear`,
				`decimals`, `dectime`, `deg2rad`, `dict2inst`, `ease`,
				`exp`, `floor`, `fmod`, `fposmod`, `funcref`,
				`hash`, `inst2dict`, `instance_from_id`, `is_inf`, `is_nan`,
				`lerp`, `linear2db`, `load`, `log`, `max`,
				`min`, `nearest_po2`, `pow`, `preload`, `print`,
				`print_stack`, `printerr`, `printraw`, `prints`, `printt`,
				`rad2deg`, `rand_range`, `rand_seed`, `randf`, `randi`,
				`randomize`, `range`, `round`, `seed`, `sign`,
				`sin`, `sinh`, `sqrt`, `stepify`, `str`,
				`str2var`, `tan`, `tanh`, `type_exist`, `typeof`,
				`var2bytes`, `var2str`, `weakref`, `yield`,
			), NameBuiltin, nil},
			{`(?<!\.)(self|false|true|PI|NAN|INF)\b`, NameBuiltinPseudo, nil},
			{Words(`(?<!\.)`, `\b`,
				`AABB`, `AcceptDialog`, `AnimatedSprite`, `AnimatedSprite3D`, `Animation`, `AnimationPlayer`, `AnimationTreePlayer`, `Area`, `Area2D`, `Array`, `AtlasTexture`, `AudioServer`, `AudioServerSW`, `AudioStream`, `AudioStreamMPC`, `AudioStreamOGGVorbis`, `AudioStreamOpus`, `AudioStreamPlayback`, `AudioStreamSpeex`, `BackBufferCopy`, `BakedLight`, `BakedLightInstance`, `BakedLightSampler`, `BaseButton`, `BitMap`, `BoneAttachment`, `bool`, `BoxContainer`, `BoxShape`, `Button`, `ButtonArray`, `ButtonGroup`, `Camera`, `Camera2D`, `CanvasItem`, `CanvasItemMaterial`, `CanvasItemShader`, `CanvasItemShaderGraph`, `CanvasLayer`, `CanvasModulate`, `CapsuleShape`, `CapsuleShape2D`, `CenterContainer`, `CheckBox`, `CheckButton`, `CircleShape2D`, `CollisionObject`, `CollisionObject2D`, `CollisionPolygon`, `CollisionPolygon2D`, `CollisionShape`, `CollisionShape2D`, `Color`, `ColorArray`, `ColorPicker`, `ColorPickerButton`, `ColorRamp`, `ConcavePolygonShape`, `ConcavePolygonShape2D`, `ConeTwistJoint`, `ConfigFile`, `ConfirmationDialog`, `Container`, `Control`, `ConvexPolygonShape`, `ConvexPolygonShape2D`, `CubeMap`, `Curve2D`, `Curve3D`, `DampedSpringJoint2D`, `Dictionary`, `DirectionalLight`, `Directory`, `EditorFileDialog`, `EditorImportPlugin`, `EditorPlugin`, `EditorScenePostImport`, `EditorScript`, `Environment`, `EventPlayer`, `EventStream`, `EventStreamChibi`, `File`, `FileDialog`, `FixedMaterial`, `float`, `Font`, `FuncRef`, `GDFunctionState`, `GDNativeClass`, `GDScript`, `Generic6DOFJoint`, `Geometry`, `GeometryInstance`, `Globals`, `GraphEdit`, `GraphNode`, `GridContainer`, `GridMap`, `GrooveJoint2D`, `HBoxContainer`, `HButtonArray`, `HingeJoint`, `HScrollBar`, `HSeparator`, `HSlider`, `HSplitContainer`, `HTTPClient`, `Image`, `ImageTexture`, `ImmediateGeometry`, `Input`, `InputDefault`, `InputEvent`, `InputEventAction`, `InputEventJoystickButton`, `InputEventJoystickMotion`, `InputEventKey`, `InputEventMouseButton`, `InputEventMouseMotion`, `InputEventScreenDrag`, `InputEventScreenTouch`, `InputMap`, `InstancePlaceholder`, `int`, `IntArray`, `InterpolatedCamera`, `IP`, `IP_Unix`, `ItemList`, `Joint`, `Joint2D`, `KinematicBody`, `KinematicBody2D`, `Label`, `LargeTexture`, `Light`, `Light2D`, `LightOccluder2D`, `LineEdit`, `LineShape2D`, `MainLoop`, `MarginContainer`, `Marshalls`, `Material`, `MaterialShader`, `MaterialShaderGraph`, `Matrix3`, `Matrix32`, `MenuButton`, `Mesh`, `MeshDataTool`, `MeshInstance`, `MeshLibrary`, `MultiMesh`, `MultiMeshInstance`, `Mutex`, `Navigation`, `Navigation2D`, `NavigationMesh`, `NavigationMeshInstance`, `NavigationPolygon`, `NavigationPolygonInstance`, `Nil`, `Node`, `Node2D`, `NodePath`, `Object`, `OccluderPolygon2D`, `OmniLight`, `OptionButton`, `OS`, `PackedDataContainer`, `PackedDataContainerRef`, `PackedScene`, `PacketPeer`, `PacketPeerStream`, `PacketPeerUDP`, `Panel`, `PanelContainer`, `ParallaxBackground`, `ParallaxLayer`, `ParticleAttractor2D`, `Particles`, `Particles2D`, `Patch9Frame`, `Path`, `Path2D`, `PathFollow`, `PathFollow2D`, `PathRemap`, `PCKPacker`, `Performance`, `PHashTranslation`, `Physics2DDirectBodyState`, `Physics2DDirectBodyStateSW`, `Physics2DDirectSpaceState`, `Physics2DServer`, `Physics2DServerSW`, `Physics2DShapeQueryParameters`, `Physics2DShapeQueryResult`, `Physics2DTestMotionResult`, `PhysicsBody`, `PhysicsBody2D`, `PhysicsDirectBodyState`, `PhysicsDirectBodyStateSW`, `PhysicsDirectSpaceState`, `PhysicsServer`, `PhysicsServerSW`, `PhysicsShapeQueryParameters`, `PhysicsShapeQueryResult`, `PinJoint`, `PinJoint2D`, `Plane`, `PlaneShape`, `Polygon2D`, `PolygonPathFinder`, `Popup`, `PopupDialog`, `PopupMenu`, `PopupPanel`, `Portal`, `Position2D`, `Position3D`, `ProgressBar`, `ProximityGroup`, `Quad`, `Quat`, `Range`, `RawArray`, `RayCast`, `RayCast2D`, `RayShape`, `RayShape2D`, `RealArray`, `Rect2`, `RectangleShape2D`, `Reference`, `ReferenceFrame`, `RegEx`, `RemoteTransform2D`, `RenderTargetTexture`, `Resource`, `ResourceImportMetadata`, `ResourceInteractiveLoader`, `ResourceLoader`, `ResourcePreloader`, `ResourceSaver`, `RichTextLabel`, `RID`, `RigidBody`, `RigidBody2D`, `Room`, `RoomBounds`, `Sample`, `SampleLibrary`, `SamplePlayer`, `SamplePlayer2D`, `SceneState`, `SceneTree`, `Script`, `ScrollBar`, `ScrollContainer`, `SegmentShape2D`, `Semaphore`, `Separator`, `Shader`, `ShaderGraph`, `ShaderMaterial`, `Shape`, `Shape2D`, `Skeleton`, `Slider`, `SliderJoint`, `SoundPlayer2D`, `SoundRoomParams`, `Spatial`, `SpatialPlayer`, `SpatialSamplePlayer`, `SpatialSound2DServer`, `SpatialSound2DServerSW`, `SpatialSoundServer`, `SpatialSoundServerSW`, `SpatialStreamPlayer`, `SphereShape`, `SpinBox`, `SplitContainer`, `SpotLight`, `Sprite`, `Sprite3D`, `SpriteBase3D`, `SpriteFrames`, `StaticBody`, `StaticBody2D`, `StreamPeer`, `StreamPeerSSL`, `StreamPeerTCP`, `StreamPlayer`, `String`, `StringArray`, `StyleBox`, `StyleBoxEmpty`, `StyleBoxFlat`, `StyleBoxImageMask`, `StyleBoxTexture`, `SurfaceTool`, `TabContainer`, `Tabs`, `TCP_Server`, `TestCube`, `TextEdit`, `Texture`, `TextureButton`, `TextureFrame`, `TextureProgress`, `Theme`, `Thread`, `TileMap`, `TileSet`, `Timer`, `ToolButton`, `TouchScreenButton`, `Transform`, `Translation`, `TranslationServer`, `Tree`, `TreeItem`, `Tween`, `UndoRedo`, `VBoxContainer`, `VButtonArray`, `Vector2`, `Vector2Array`, `Vector3`, `Vector3Array`, `VehicleBody`, `VehicleWheel`, `VideoPlayer`, `VideoStream`, `VideoStreamTheora`, `Viewport`, `ViewportSprite`, `VisibilityEnabler`, `VisibilityEnabler2D`, `VisibilityNotifier`, `VisibilityNotifier2D`, `VisualInstance`, `VisualServer`, `VScrollBar`, `VSeparator`, `VSlider`, `VSplitContainer`, `WeakRef`, `WindowDialog`, `World`, `World2D`, `WorldEnvironment`, `XMLParser`, `YSort`), NameException, nil},
		},
		"numbers": {
			{`(\d+\.\d*|\d*\.\d+)([eE][+-]?[0-9]+)?j?`, LiteralNumberFloat, nil},
			{`\d+[eE][+-]?[0-9]+j?`, LiteralNumberFloat, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumberHex, nil},
			{`\d+j?`, LiteralNumberInteger, nil},
		},
		"name": {
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"funcname": {
			{`[a-zA-Z_]\w*`, NameFunction, Pop(1)},
			Default(Pop(1)),
		},
		"classname": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
		},
		"stringescape": {
			{`\\([\\abfnrtv"\']|\n|N\{.*?\}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|x[a-fA-F0-9]{2}|[0-7]{1,3})`, LiteralStringEscape, nil},
		},
		"strings-single": {
			{`%(\(\w+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsux%]`, LiteralStringInterpol, nil},
			{`[^\\\'"%\n]+`, LiteralStringSingle, nil},
			{`[\'"\\]`, LiteralStringSingle, nil},
			{`%`, LiteralStringSingle, nil},
		},
		"strings-double": {
			{`%(\(\w+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsux%]`, LiteralStringInterpol, nil},
			{`[^\\\'"%\n]+`, LiteralStringDouble, nil},
			{`[\'"\\]`, LiteralStringDouble, nil},
			{`%`, LiteralStringDouble, nil},
		},
		"dqs": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`\\\\|\\"|\\\n`, LiteralStringEscape, nil},
			Include("strings-double"),
		},
		"sqs": {
			{`'`, LiteralStringSingle, Pop(1)},
			{`\\\\|\\'|\\\n`, LiteralStringEscape, nil},
			Include("strings-single"),
		},
		"tdqs": {
			{`"""`, LiteralStringDouble, Pop(1)},
			Include("strings-double"),
			{`\n`, LiteralStringDouble, nil},
		},
		"tsqs": {
			{`'''`, LiteralStringSingle, Pop(1)},
			Include("strings-single"),
			{`\n`, LiteralStringSingle, nil},
		},
	}
}
